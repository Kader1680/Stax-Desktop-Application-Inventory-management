package domain

import "time"

type Category struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Product struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	SKU           string    `json:"sku"`
	Barcode       string    `json:"barcode,omitempty"`
	CategoryID    *int64    `json:"categoryId,omitempty"`
	CategoryName  string    `json:"categoryName,omitempty"` // populated on read via JOIN
	Unit          Unit      `json:"unit"`
	PurchasePrice Money     `json:"purchasePrice"`
	SellingPrice  Money     `json:"sellingPrice"`
	StockQty      Quantity  `json:"stockQty"`
	MinStockQty   Quantity  `json:"minStockQty"`
	IsActive      bool      `json:"isActive"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// IsLowStock reports whether current stock is at or below the minimum threshold.
func (p Product) IsLowStock() bool {
	return !p.StockQty.LessThan(ZeroQuantity()) && !p.MinStockQty.Decimal.IsZero() &&
		(p.StockQty.LessThan(p.MinStockQty) || p.StockQty.Decimal.Equal(p.MinStockQty.Decimal))
}

type MovementType string

const (
	MovementIn         MovementType = "in"
	MovementOut        MovementType = "out"
	MovementAdjustment MovementType = "adjustment"
	MovementSale       MovementType = "sale"
	MovementPurchase   MovementType = "purchase"
)

type StockMovement struct {
	ID            int64        `json:"id"`
	ProductID     int64        `json:"productId"`
	Type          MovementType `json:"type"`
	Qty           Quantity     `json:"qty"` // signed: positive = stock increased, negative = decreased
	Reason        string       `json:"reason,omitempty"`
	ReferenceID   *int64       `json:"referenceId,omitempty"`
	ReferenceType string       `json:"referenceType,omitempty"`
	CreatedAt     time.Time    `json:"createdAt"`
}

// ProductInput is the payload accepted from the UI when creating/editing a product.
// Kept separate from Product so the frontend sends plain strings/floats we validate here.
type ProductInput struct {
	Name          string  `json:"name"`
	SKU           string  `json:"sku"`
	Barcode       string  `json:"barcode"`
	CategoryID    *int64  `json:"categoryId"`
	Unit          string  `json:"unit"`
	PurchasePrice float64 `json:"purchasePrice"`
	SellingPrice  float64 `json:"sellingPrice"`
	StockQty      string  `json:"stockQty"`
	MinStockQty   string  `json:"minStockQty"`
}
