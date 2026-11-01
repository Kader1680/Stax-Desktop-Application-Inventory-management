package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"stax-desktop/internal/domain"
	"stax-desktop/internal/repository"
)

type InventoryService struct {
	products  repository.ProductRepository
	movements repository.StockMovementRepository
	tx        repository.TxManager
}

func NewInventoryService(products repository.ProductRepository, movements repository.StockMovementRepository, tx repository.TxManager) *InventoryService {
	return &InventoryService{products: products, movements: movements, tx: tx}
}

// StockIn increases stock for a product (e.g. manual restock outside a purchase order).
func (s *InventoryService) StockIn(ctx context.Context, productID int64, qtyStr, reason string) error {
	qty, err := s.parsePositiveQty(ctx, productID, qtyStr)
	if err != nil {
		return err
	}
	return s.applyMovement(ctx, productID, domain.MovementIn, qty, reason, nil, "")
}

// StockOut decreases stock (e.g. damaged goods, internal use) — not a sale.
func (s *InventoryService) StockOut(ctx context.Context, productID int64, qtyStr, reason string) error {
	qty, err := s.parsePositiveQty(ctx, productID, qtyStr)
	if err != nil {
		return err
	}
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("%w: une raison est requise pour une sortie de stock", domain.ErrInvalidInput)
	}
	negative := domain.ZeroQuantity().Sub(qty)
	return s.applyMovement(ctx, productID, domain.MovementOut, negative, reason, nil, "")
}

// Adjustment sets an absolute correction (delta can be positive or negative),
// e.g. after a physical stock count. A reason is mandatory for audit purposes.
func (s *InventoryService) Adjustment(ctx context.Context, productID int64, deltaStr, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return fmt.Errorf("%w: une raison est requise pour un ajustement", domain.ErrInvalidInput)
	}
	delta, err := domain.NewQuantity(deltaStr)
	if err != nil {
		return fmt.Errorf("%w: quantité d'ajustement invalide", domain.ErrInvalidInput)
	}
	return s.applyMovement(ctx, productID, domain.MovementAdjustment, delta, reason, nil, "")
}

func (s *InventoryService) History(ctx context.Context, productID int64, limit int) ([]domain.StockMovement, error) {
	return s.movements.ListByProduct(ctx, productID, limit)
}

func (s *InventoryService) RecentMovements(ctx context.Context, limit int) ([]domain.StockMovement, error) {
	return s.movements.ListRecent(ctx, limit)
}

// ApplyStockChangeInTx is used by other services (POS, Purchases) that already
// hold an open transaction (e.g. creating a sale + decrementing stock atomically).
// signedQty: negative for sales/stock-out, positive for purchases/stock-in.
func (s *InventoryService) ApplyStockChangeInTx(ctx context.Context, tx *sql.Tx, productID int64, mType domain.MovementType, signedQty domain.Quantity, reason string, refID *int64, refType string) error {
	if err := s.products.AdjustStock(ctx, tx, productID, signedQty); err != nil {
		return err
	}
	_, err := s.movements.Create(ctx, tx, &domain.StockMovement{
		ProductID:     productID,
		Type:          mType,
		Qty:           signedQty,
		Reason:        reason,
		ReferenceID:   refID,
		ReferenceType: refType,
	})
	return err
}

func (s *InventoryService) parsePositiveQty(ctx context.Context, productID int64, qtyStr string) (domain.Quantity, error) {
	qty, err := domain.NewQuantity(qtyStr)
	if err != nil {
		return domain.Quantity{}, fmt.Errorf("%w: quantité invalide", domain.ErrInvalidInput)
	}
	if qty.IsNegative() || qty.Decimal.IsZero() {
		return domain.Quantity{}, fmt.Errorf("%w: la quantité doit être positive", domain.ErrInvalidInput)
	}
	product, err := s.products.GetByID(ctx, productID)
	if err != nil {
		return domain.Quantity{}, err
	}
	if err := qty.Validate(product.Unit); err != nil {
		return domain.Quantity{}, fmt.Errorf("%w: %s", domain.ErrInvalidInput, err.Error())
	}
	return qty, nil
}

// applyMovement is the single transactional chokepoint: adjust stock + write
// the audit row atomically, used by inventory ops, sales and purchases.
func (s *InventoryService) applyMovement(ctx context.Context, productID int64, mType domain.MovementType, signedQty domain.Quantity, reason string, refID *int64, refType string) error {
	return s.tx.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.products.AdjustStock(ctx, tx, productID, signedQty); err != nil {
			return err
		}
		_, err := s.movements.Create(ctx, tx, &domain.StockMovement{
			ProductID:     productID,
			Type:          mType,
			Qty:           signedQty,
			Reason:        reason,
			ReferenceID:   refID,
			ReferenceType: refType,
		})
		return err
	})
}
