package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"stax-desktop/internal/domain"
	"stax-desktop/internal/repository"
)

type ProductService struct {
	products  repository.ProductRepository
	movements repository.StockMovementRepository
	tx        repository.TxManager
}

func NewProductService(products repository.ProductRepository, movements repository.StockMovementRepository, tx repository.TxManager) *ProductService {
	return &ProductService{products: products, movements: movements, tx: tx}
}

// Create validates the input and creates a new product, recording an initial
// "in" stock movement if an opening stock quantity was provided.
func (s *ProductService) Create(ctx context.Context, in domain.ProductInput) (*domain.Product, error) {
	p, err := s.validateAndBuild(ctx, 0, in)
	if err != nil {
		return nil, err
	}

	id, err := s.products.Create(ctx, p)
	if err != nil {
		return nil, err
	}
	p.ID = id

	// Record the opening stock as an auditable movement, inside its own tx.
	if !p.StockQty.Decimal.IsZero() {
		err = s.tx.WithTx(ctx, func(tx *sql.Tx) error {
			_, err := s.movements.Create(ctx, tx, &domain.StockMovement{
				ProductID: id,
				Type:      domain.MovementIn,
				Qty:       p.StockQty,
				Reason:    "Stock initial à la création du produit",
			})
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("produit créé mais échec de l'enregistrement du mouvement de stock: %w", err)
		}
	}
	return p, nil
}

// Update validates and persists changes to an existing product.
// Stock quantity is NOT editable here directly — it must go through
// InventoryService (StockIn/StockOut/Adjustment) so every change is audited.
func (s *ProductService) Update(ctx context.Context, id int64, in domain.ProductInput) (*domain.Product, error) {
	existing, err := s.products.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	p, err := s.validateAndBuild(ctx, id, in)
	if err != nil {
		return nil, err
	}
	p.StockQty = existing.StockQty // preserve; stock changes go through InventoryService

	if err := s.products.Update(ctx, p); err != nil {
		return nil, err
	}
	p.ID = id
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id int64) error {
	return s.products.SoftDelete(ctx, id)
}

func (s *ProductService) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	return s.products.GetByID(ctx, id)
}

func (s *ProductService) FindByBarcode(ctx context.Context, barcode string) (*domain.Product, error) {
	return s.products.FindByBarcode(ctx, barcode)
}

func (s *ProductService) Search(ctx context.Context, query string, categoryID *int64) ([]domain.Product, error) {
	return s.products.Search(ctx, strings.TrimSpace(query), categoryID, true)
}

func (s *ProductService) LowStock(ctx context.Context) ([]domain.Product, error) {
	return s.products.LowStock(ctx)
}

// validateAndBuild applies all business validation rules and returns a Product
// ready to persist. excludeID is used when updating (0 means "new product").
func (s *ProductService) validateAndBuild(ctx context.Context, excludeID int64, in domain.ProductInput) (*domain.Product, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: le nom du produit est requis", domain.ErrInvalidInput)
	}
	sku := strings.TrimSpace(in.SKU)
	if sku == "" {
		return nil, fmt.Errorf("%w: le SKU est requis", domain.ErrInvalidInput)
	}

	unit := domain.Unit(in.Unit)
	if !unit.IsValid() {
		return nil, fmt.Errorf("%w: unité '%s' non supportée", domain.ErrInvalidUnit, in.Unit)
	}

	if in.SellingPrice < 0 || in.PurchasePrice < 0 {
		return nil, fmt.Errorf("%w: les prix ne peuvent pas être négatifs", domain.ErrInvalidInput)
	}

	// Duplicate SKU check (excluding self on update).
	if existing, err := s.products.FindBySKU(ctx, sku); err == nil && existing.ID != excludeID {
		return nil, domain.ErrDuplicateSKU
	}

	barcode := strings.TrimSpace(in.Barcode)
	if barcode != "" {
		if existing, err := s.products.FindByBarcode(ctx, barcode); err == nil && existing.ID != excludeID {
			return nil, domain.ErrDuplicateBarcode
		}
	}

	stockQtyStr := in.StockQty
	if stockQtyStr == "" {
		stockQtyStr = "0"
	}
	stockQty, err := domain.NewQuantity(stockQtyStr)
	if err != nil {
		return nil, fmt.Errorf("%w: quantité en stock invalide", domain.ErrInvalidInput)
	}
	if err := stockQty.Validate(unit); err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrInvalidInput, err.Error())
	}

	minQtyStr := in.MinStockQty
	if minQtyStr == "" {
		minQtyStr = "0"
	}
	minQty, err := domain.NewQuantity(minQtyStr)
	if err != nil {
		return nil, fmt.Errorf("%w: seuil minimum invalide", domain.ErrInvalidInput)
	}

	return &domain.Product{
		Name:          name,
		SKU:           sku,
		Barcode:       barcode,
		CategoryID:    in.CategoryID,
		Unit:          unit,
		PurchasePrice: domain.NewMoneyFromFloat(in.PurchasePrice),
		SellingPrice:  domain.NewMoneyFromFloat(in.SellingPrice),
		StockQty:      stockQty,
		MinStockQty:   minQty,
		IsActive:      true,
	}, nil
}
