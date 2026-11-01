package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"stax-desktop/internal/domain"
)

type ProductRepo struct {
	DB *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{DB: db}
}

const productSelectCols = `
	p.id, p.name, p.sku, COALESCE(p.barcode,''), p.category_id,
	COALESCE(c.name,''), p.unit, p.purchase_price, p.selling_price,
	p.stock_qty, p.min_stock_qty, p.is_active, p.created_at, p.updated_at
`

func (r *ProductRepo) scanProduct(row interface{ Scan(...any) error }) (*domain.Product, error) {
	var p domain.Product
	var categoryID sql.NullInt64
	var stockQtyStr, minQtyStr, unit string
	var createdAt, updatedAt string

	err := row.Scan(
		&p.ID, &p.Name, &p.SKU, &p.Barcode, &categoryID,
		&p.CategoryName, &unit, &p.PurchasePrice, &p.SellingPrice,
		&stockQtyStr, &minQtyStr, &p.IsActive, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	if categoryID.Valid {
		p.CategoryID = &categoryID.Int64
	}
	p.Unit = domain.Unit(unit)
	if p.StockQty, err = domain.NewQuantity(stockQtyStr); err != nil {
		return nil, err
	}
	if p.MinStockQty, err = domain.NewQuantity(minQtyStr); err != nil {
		return nil, err
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updatedAt)
	return &p, nil
}

func (r *ProductRepo) Create(ctx context.Context, p *domain.Product) (int64, error) {
	res, err := r.DB.ExecContext(ctx, `
		INSERT INTO products (name, sku, barcode, category_id, unit, purchase_price,
			selling_price, stock_qty, min_stock_qty, is_active)
		VALUES (?, ?, NULLIF(?, ''), ?, ?, ?, ?, ?, ?, 1)`,
		p.Name, p.SKU, p.Barcode, p.CategoryID, string(p.Unit), int64(p.PurchasePrice),
		int64(p.SellingPrice), p.StockQty.String(), p.MinStockQty.String(),
	)
	if err != nil {
		return 0, mapUniqueConstraint(err)
	}
	return res.LastInsertId()
}

func (r *ProductRepo) Update(ctx context.Context, p *domain.Product) error {
	_, err := r.DB.ExecContext(ctx, `
		UPDATE products SET name=?, sku=?, barcode=NULLIF(?, ''), category_id=?, unit=?,
			purchase_price=?, selling_price=?, min_stock_qty=?, updated_at=strftime('%Y-%m-%dT%H:%M:%fZ','now')
		WHERE id=?`,
		p.Name, p.SKU, p.Barcode, p.CategoryID, string(p.Unit),
		int64(p.PurchasePrice), int64(p.SellingPrice), p.MinStockQty.String(), p.ID,
	)
	if err != nil {
		return mapUniqueConstraint(err)
	}
	return nil
}

func (r *ProductRepo) SoftDelete(ctx context.Context, id int64) error {
	res, err := r.DB.ExecContext(ctx, `UPDATE products SET is_active=0 WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT `+productSelectCols+`
		FROM products p LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.id = ?`, id)
	p, err := r.scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *ProductRepo) FindBySKU(ctx context.Context, sku string) (*domain.Product, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT `+productSelectCols+`
		FROM products p LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.sku = ?`, sku)
	p, err := r.scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *ProductRepo) FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	row := r.DB.QueryRowContext(ctx, `
		SELECT `+productSelectCols+`
		FROM products p LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.barcode = ?`, barcode)
	p, err := r.scanProduct(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *ProductRepo) Search(ctx context.Context, query string, categoryID *int64, onlyActive bool) ([]domain.Product, error) {
	sb := strings.Builder{}
	sb.WriteString(`SELECT ` + productSelectCols + ` FROM products p LEFT JOIN categories c ON c.id = p.category_id WHERE 1=1`)
	args := []any{}

	if query != "" {
		sb.WriteString(` AND (p.name LIKE ? OR p.sku LIKE ? OR p.barcode LIKE ?)`)
		like := "%" + query + "%"
		args = append(args, like, like, like)
	}
	if categoryID != nil {
		sb.WriteString(` AND p.category_id = ?`)
		args = append(args, *categoryID)
	}
	if onlyActive {
		sb.WriteString(` AND p.is_active = 1`)
	}
	sb.WriteString(` ORDER BY p.name ASC LIMIT 500`)

	rows, err := r.DB.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Product
	for rows.Next() {
		p, err := r.scanProduct(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, rows.Err()
}

func (r *ProductRepo) LowStock(ctx context.Context) ([]domain.Product, error) {
	// Compare as real numbers since SQLite stores qty as TEXT decimal strings.
	rows, err := r.DB.QueryContext(ctx, `
		SELECT `+productSelectCols+`
		FROM products p LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.is_active = 1 AND CAST(p.stock_qty AS REAL) <= CAST(p.min_stock_qty AS REAL)
			AND CAST(p.min_stock_qty AS REAL) > 0
		ORDER BY p.name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Product
	for rows.Next() {
		p, err := r.scanProduct(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, rows.Err()
}

// AdjustStock updates stock_qty by delta (positive or negative) using a decimal-safe
// read-modify-write inside the caller's transaction. Returns domain.ErrInsufficientStock
// if the resulting stock would go negative.
func (r *ProductRepo) AdjustStock(ctx context.Context, tx *sql.Tx, productID int64, delta domain.Quantity) error {
	var currentStr string
	err := tx.QueryRowContext(ctx, `SELECT stock_qty FROM products WHERE id = ?`, productID).Scan(&currentStr)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	if err != nil {
		return err
	}
	current, err := domain.NewQuantity(currentStr)
	if err != nil {
		return err
	}
	newQty := current.Add(delta)
	if newQty.IsNegative() {
		return domain.ErrInsufficientStock
	}
	_, err = tx.ExecContext(ctx, `
		UPDATE products SET stock_qty = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ?`,
		newQty.String(), productID,
	)
	return err
}

func mapUniqueConstraint(err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "products.sku"):
		return domain.ErrDuplicateSKU
	case strings.Contains(msg, "products.barcode"):
		return domain.ErrDuplicateBarcode
	default:
		return fmt.Errorf("erreur base de données: %w", err)
	}
}
