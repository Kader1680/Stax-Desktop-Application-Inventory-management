package repository

import (
	"context"
	"database/sql"

	"stax-desktop/internal/domain"
)

// ProductRepository persists and queries products.
type ProductRepository interface {
	Create(ctx context.Context, p *domain.Product) (int64, error)
	Update(ctx context.Context, p *domain.Product) error
	SoftDelete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (*domain.Product, error)
	FindBySKU(ctx context.Context, sku string) (*domain.Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error)
	Search(ctx context.Context, query string, categoryID *int64, onlyActive bool) ([]domain.Product, error)
	LowStock(ctx context.Context) ([]domain.Product, error)
	// AdjustStock changes stock_qty by delta (can be negative) within the given tx.
	AdjustStock(ctx context.Context, tx *sql.Tx, productID int64, delta domain.Quantity) error
}

// StockMovementRepository records and queries the inventory audit trail.
type StockMovementRepository interface {
	Create(ctx context.Context, tx *sql.Tx, m *domain.StockMovement) (int64, error)
	ListByProduct(ctx context.Context, productID int64, limit int) ([]domain.StockMovement, error)
	ListRecent(ctx context.Context, limit int) ([]domain.StockMovement, error)
}

// CategoryRepository persists and queries categories.
type CategoryRepository interface {
	Create(ctx context.Context, c *domain.Category) (int64, error)
	List(ctx context.Context) ([]domain.Category, error)
	Delete(ctx context.Context, id int64) error
}

// TxManager runs a function inside a database transaction, committing on success
// and rolling back on error or panic. Used by services for multi-step operations.
type TxManager interface {
	WithTx(ctx context.Context, fn func(tx *sql.Tx) error) error
}
