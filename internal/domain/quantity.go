package domain

import (
	"fmt"

	"github.com/shopspring/decimal"
)

// Unit is a product's unit of measure.
type Unit string

const (
	UnitPiece     Unit = "piece"
	UnitKilogram  Unit = "kg"
	UnitGram      Unit = "g"
	UnitLiter     Unit = "l"
)

func (u Unit) IsValid() bool {
	switch u {
	case UnitPiece, UnitKilogram, UnitGram, UnitLiter:
		return true
	}
	return false
}

// ValidUnits lists all supported units, used for frontend dropdowns and validation.
var ValidUnits = []Unit{UnitPiece, UnitKilogram, UnitGram, UnitLiter}

// Quantity wraps decimal.Decimal so stock/sale quantities never lose precision
// (e.g. 0.750 kg, 1.5 l). Stored in SQLite as TEXT.
type Quantity struct {
	decimal.Decimal
}

func NewQuantity(s string) (Quantity, error) {
	d, err := decimal.NewFromString(s)
	if err != nil {
		return Quantity{}, fmt.Errorf("invalid quantity %q: %w", s, err)
	}
	return Quantity{d}, nil
}

func QuantityFromFloat(f float64) Quantity {
	return Quantity{decimal.NewFromFloat(f)}
}

func ZeroQuantity() Quantity {
	return Quantity{decimal.Zero}
}

func (q Quantity) Add(other Quantity) Quantity {
	return Quantity{q.Decimal.Add(other.Decimal)}
}

func (q Quantity) Sub(other Quantity) Quantity {
	return Quantity{q.Decimal.Sub(other.Decimal)}
}

func (q Quantity) IsNegative() bool {
	return q.Decimal.IsNegative()
}

func (q Quantity) LessThan(other Quantity) bool {
	return q.Decimal.LessThan(other.Decimal)
}

// RequiresWholeNumber reports whether a unit must use integer quantities (pieces).
func RequiresWholeNumber(u Unit) bool {
	return u == UnitPiece
}

// Validate checks a quantity is consistent with its unit's rules (e.g. pieces can't be fractional).
func (q Quantity) Validate(u Unit) error {
	if RequiresWholeNumber(u) && !q.Decimal.Equal(q.Decimal.Truncate(0)) {
		return fmt.Errorf("la quantité pour l'unité 'pièce' doit être un nombre entier")
	}
	if q.IsNegative() {
		return fmt.Errorf("la quantité ne peut pas être négative")
	}
	return nil
}
