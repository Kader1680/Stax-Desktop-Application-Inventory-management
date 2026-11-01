# STAX Desktop

Offline-first wholesale/retail inventory & POS desktop app (Go + Wails + React/TS + SQLite).

## Status

Scaffolded and implemented so far:
- Domain layer: `Money`, `Quantity` (decimal/unit-safe), `Product`, `Category`, `StockMovement`.
- Database layer: embedded SQL migrations (`internal/database/migrations`), WAL mode, FK enforcement, `VACUUM INTO` backup.
- Repository layer: `ProductRepo`, `CategoryRepo`, `StockMovementRepo`, `TxManager` (sqlite).
- Service layer: `ProductService` (validation, CRUD, duplicate SKU/barcode checks), `InventoryService` (stock in/out/adjustment, fully transactional, audit trail).
- Wails bindings: `internal/app/App` exposing Products + Inventory to the frontend.
- French i18n seed (`frontend/src/i18n/fr.json`) — structure supports adding `ar.json` + RTL later without rearchitecting.

Not yet implemented (next steps, in order): POS/Sales service + UI, Purchases, Customers/Suppliers, Reports, Dashboard, Demo mode, React pages, backup/restore UI.

## Why this wasn't compiled here

This sandbox has no Go toolchain and can't produce a Windows GUI binary regardless. **I'd strongly recommend continuing this in Claude Code on your own machine** — it can run `go build`, `wails dev`, and catch compile errors as we go, which matters a lot for a multi-layer Go project like this.

## Local setup (Windows dev machine)

1. Install Go 1.22+, Node 18+, and Wails CLI:
   ```
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```
2. From the project root:
   ```
   go mod tidy
   cd frontend && npm install && cd ..
   wails dev
   ```
3. `go mod tidy` will pull in `github.com/golang-migrate/migrate/v4` — add it to `go.mod` if not already present (it's used in `internal/database/db.go` but not yet listed in `go.mod` here; `go mod tidy` fixes this automatically).

## Project structure

```
internal/
  domain/        — models, Money, Quantity, sentinel errors
  database/      — SQLite connection + embedded migrations
  repository/    — interfaces + sqlite/ implementations
  service/       — business logic (ProductService, InventoryService, ...)
  app/           — Wails-bound App struct (thin controller layer)
frontend/src/
  i18n/          — fr.json now, ar.json later for RTL
main.go
wails.json
```

## Design notes carried through the code

- **Money** = integer cents (`domain.Money`), never float, to avoid rounding drift.
- **Quantity** = decimal string via `shopspring/decimal`, supports fractional kg/L, blocks fractional "piece" units.
- **Stock is never edited directly.** All stock changes go through `InventoryService`, which writes a `stock_movements` row in the *same transaction* as the stock update — so stock levels are always reconstructible/auditable.
- **Soft delete** on products (`is_active`), so sale/purchase history stays valid after a product is "deleted".
- Errors are sentinel-based in `domain/errors.go`, translated to French at the service layer, and the `app` layer guarantees no raw Go/SQL error ever reaches the UI.

## Next module: POS / Sales

The sale flow will reuse `InventoryService.ApplyStockChangeInTx` inside a `SaleService.Confirm(...)` transaction: create `sales` + `sale_items` rows, decrement stock per line, write `stock_movements(type=sale, reference_id=sale.id)`, generate `invoice_no`, all atomically. I'd build and verify that next, the same way as Products/Inventory above — happy to continue as soon as you've confirmed this module builds cleanly on your machine.
