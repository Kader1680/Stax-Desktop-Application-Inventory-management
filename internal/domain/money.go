package domain

import "fmt"

// Money represents an amount in the smallest currency unit (cents/centimes).
// Stored and transmitted as int64 to avoid floating point rounding errors.
type Money int64

// NewMoneyFromFloat converts a decimal amount (e.g. 19.99) to Money.
// Rounds to the nearest cent.
func NewMoneyFromFloat(amount float64) Money {
	return Money(int64(amount*100 + 0.5))
}

// Float returns the money as a float64 (for display only, never for calculations).
func (m Money) Float() float64 {
	return float64(m) / 100.0
}

func (m Money) Add(other Money) Money      { return m + other }
func (m Money) Sub(other Money) Money      { return m - other }
func (m Money) Mul(times int64) Money      { return m * Money(times) }
func (m Money) IsNegative() bool           { return m < 0 }
func (m Money) String() string             { return fmt.Sprintf("%.2f", m.Float()) }

// ApplyPercentDiscount returns the discount amount (not the remainder) for a given percent (0-100).
func (m Money) ApplyPercentDiscount(percent float64) Money {
	return Money(int64(float64(m) * percent / 100.0))
}
