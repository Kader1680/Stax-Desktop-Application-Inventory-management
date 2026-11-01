package app

import (
	"context"
	"database/sql"
	"errors"

	"stax-desktop/internal/domain"
	"stax-desktop/internal/repository/sqlite"
	"stax-desktop/internal/service"
)

// App is bound to the Wails runtime. Every exported method becomes callable
// from the React frontend via the generated bindings. Keep this layer thin:
// no business logic here, just translating calls <-> services and errors
// into French user-facing messages.
type App struct {
	ctx context.Context

	products  *service.ProductService
	inventory *service.InventoryService
}

func NewApp(db *sql.DB) *App {
	productRepo := sqlite.NewProductRepo(db)
	movementRepo := sqlite.NewStockMovementRepo(db)
	txManager := sqlite.NewTxManager(db)

	return &App{
		products:  service.NewProductService(productRepo, movementRepo, txManager),
		inventory: service.NewInventoryService(productRepo, movementRepo, txManager),
	}
}

// Startup is called by Wails when the app starts; ctx is required for all
// further runtime calls (e.g. dialogs, events) if/when we need them.
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// --- Products ---

func (a *App) CreateProduct(in domain.ProductInput) (*domain.Product, error) {
	p, err := a.products.Create(a.ctx, in)
	return p, friendlyError(err)
}

func (a *App) UpdateProduct(id int64, in domain.ProductInput) (*domain.Product, error) {
	p, err := a.products.Update(a.ctx, id, in)
	return p, friendlyError(err)
}

func (a *App) DeleteProduct(id int64) error {
	return friendlyError(a.products.Delete(a.ctx, id))
}

func (a *App) GetProduct(id int64) (*domain.Product, error) {
	p, err := a.products.GetByID(a.ctx, id)
	return p, friendlyError(err)
}

func (a *App) FindProductByBarcode(barcode string) (*domain.Product, error) {
	p, err := a.products.FindByBarcode(a.ctx, barcode)
	return p, friendlyError(err)
}

func (a *App) SearchProducts(query string, categoryID *int64) ([]domain.Product, error) {
	list, err := a.products.Search(a.ctx, query, categoryID)
	return list, friendlyError(err)
}

func (a *App) LowStockProducts() ([]domain.Product, error) {
	list, err := a.products.LowStock(a.ctx)
	return list, friendlyError(err)
}

// --- Inventory ---

func (a *App) StockIn(productID int64, qty string, reason string) error {
	return friendlyError(a.inventory.StockIn(a.ctx, productID, qty, reason))
}

func (a *App) StockOut(productID int64, qty string, reason string) error {
	return friendlyError(a.inventory.StockOut(a.ctx, productID, qty, reason))
}

func (a *App) StockAdjustment(productID int64, delta string, reason string) error {
	return friendlyError(a.inventory.Adjustment(a.ctx, productID, delta, reason))
}

func (a *App) ProductStockHistory(productID int64, limit int) ([]domain.StockMovement, error) {
	list, err := a.inventory.History(a.ctx, productID, limit)
	return list, friendlyError(err)
}

func (a *App) RecentStockMovements(limit int) ([]domain.StockMovement, error) {
	list, err := a.inventory.RecentMovements(a.ctx, limit)
	return list, friendlyError(err)
}

// friendlyError passes through nil and already-French domain errors as-is;
// anything unexpected is wrapped so we never leak raw Go/SQL errors to the UI.
func friendlyError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, domain.ErrNotFound) ||
		errors.Is(err, domain.ErrDuplicateSKU) ||
		errors.Is(err, domain.ErrDuplicateBarcode) ||
		errors.Is(err, domain.ErrInvalidInput) ||
		errors.Is(err, domain.ErrInsufficientStock) ||
		errors.Is(err, domain.ErrInvalidUnit) {
		return err
	}
	return errors.New("une erreur inattendue est survenue: " + err.Error())
}
