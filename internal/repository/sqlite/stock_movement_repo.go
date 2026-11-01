package sqlite

import (
	"context"
	"database/sql"
	"time"

	"stax-desktop/internal/domain"
)

type StockMovementRepo struct {
	DB *sql.DB
}

func NewStockMovementRepo(db *sql.DB) *StockMovementRepo {
	return &StockMovementRepo{DB: db}
}

// Create inserts a movement row within the given transaction (movements always
// accompany a stock change, so they share the same tx as AdjustStock).
func (r *StockMovementRepo) Create(ctx context.Context, tx *sql.Tx, m *domain.StockMovement) (int64, error) {
	res, err := tx.ExecContext(ctx, `
		INSERT INTO stock_movements (product_id, type, qty, reason, reference_id, reference_type)
		VALUES (?, ?, ?, ?, ?, ?)`,
		m.ProductID, string(m.Type), m.Qty.String(), m.Reason, m.ReferenceID, m.ReferenceType,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *StockMovementRepo) scanMovement(row interface{ Scan(...any) error }) (*domain.StockMovement, error) {
	var m domain.StockMovement
	var qtyStr, typ, createdAt string
	var refID sql.NullInt64
	var refType sql.NullString
	var reason sql.NullString

	err := row.Scan(&m.ID, &m.ProductID, &typ, &qtyStr, &reason, &refID, &refType, &createdAt)
	if err != nil {
		return nil, err
	}
	m.Type = domain.MovementType(typ)
	if m.Qty, err = domain.NewQuantity(qtyStr); err != nil {
		return nil, err
	}
	m.Reason = reason.String
	if refID.Valid {
		m.ReferenceID = &refID.Int64
	}
	m.ReferenceType = refType.String
	m.CreatedAt, _ = time.Parse(time.RFC3339Nano, createdAt)
	return &m, nil
}

func (r *StockMovementRepo) ListByProduct(ctx context.Context, productID int64, limit int) ([]domain.StockMovement, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, product_id, type, qty, reason, reference_id, reference_type, created_at
		FROM stock_movements WHERE product_id = ? ORDER BY created_at DESC, id DESC LIMIT ?`,
		productID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.collect(rows)
}

func (r *StockMovementRepo) ListRecent(ctx context.Context, limit int) ([]domain.StockMovement, error) {
	if limit <= 0 {
		limit = 50
	}
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, product_id, type, qty, reason, reference_id, reference_type, created_at
		FROM stock_movements ORDER BY created_at DESC, id DESC LIMIT ?`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return r.collect(rows)
}

func (r *StockMovementRepo) collect(rows *sql.Rows) ([]domain.StockMovement, error) {
	var out []domain.StockMovement
	for rows.Next() {
		m, err := r.scanMovement(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *m)
	}
	return out, rows.Err()
}
