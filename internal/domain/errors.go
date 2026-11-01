package domain

import "errors"

// Sentinel errors. Services wrap these; App layer maps them to French user-facing messages.
var (
	ErrNotFound          = errors.New("ressource introuvable")
	ErrDuplicateSKU      = errors.New("ce SKU existe déjà")
	ErrDuplicateBarcode  = errors.New("ce code-barres existe déjà")
	ErrInvalidInput      = errors.New("données invalides")
	ErrInsufficientStock = errors.New("stock insuffisant")
	ErrInvalidUnit       = errors.New("unité invalide")
)
