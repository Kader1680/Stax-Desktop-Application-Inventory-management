// Thin wrapper around the Wails-generated bindings (frontend/wailsjs/go/app/App.js,
// created automatically by `wails dev`/`wails build` from internal/app/app.go).
// Pages import from here instead of the generated path directly, so if the Go
// package/struct name ever changes, only this file needs updating.
import * as Generated from "../../wailsjs/go/app/App";
import type { Product, ProductInput, StockMovement } from "./types";

// Wails auto-generates its own TS types from the Go structs (wailsjs/go/models.ts),
// typing string-backed fields (unit, movement type) as plain `string` rather than
// our narrower literal unions. The runtime values are always valid (the Go side
// validates them), so we cast once here at the boundary instead of widening our
// app-side types everywhere they're used.
function as<T>(value: unknown): T {
  return value as T;
}

export const api = {
  createProduct: async (input: ProductInput): Promise<Product> =>
    as<Product>(await Generated.CreateProduct(input as any)),
  updateProduct: async (id: number, input: ProductInput): Promise<Product> =>
    as<Product>(await Generated.UpdateProduct(id, input as any)),
  deleteProduct: (id: number): Promise<void> => Generated.DeleteProduct(id),
  getProduct: async (id: number): Promise<Product> =>
    as<Product>(await Generated.GetProduct(id)),
  findProductByBarcode: async (barcode: string): Promise<Product> =>
    as<Product>(await Generated.FindProductByBarcode(barcode)),
  searchProducts: async (query: string, categoryId: number | null): Promise<Product[]> =>
    as<Product[]>(await Generated.SearchProducts(query, categoryId as any)),
  lowStockProducts: async (): Promise<Product[]> =>
    as<Product[]>(await Generated.LowStockProducts()),

  stockIn: (productId: number, qty: string, reason: string): Promise<void> =>
    Generated.StockIn(productId, qty, reason),
  stockOut: (productId: number, qty: string, reason: string): Promise<void> =>
    Generated.StockOut(productId, qty, reason),
  stockAdjustment: (productId: number, delta: string, reason: string): Promise<void> =>
    Generated.StockAdjustment(productId, delta, reason),
  productStockHistory: async (productId: number, limit: number): Promise<StockMovement[]> =>
    as<StockMovement[]>(await Generated.ProductStockHistory(productId, limit)),
  recentStockMovements: async (limit: number): Promise<StockMovement[]> =>
    as<StockMovement[]>(await Generated.RecentStockMovements(limit)),
};

/** Extracts a displayable message from an error thrown by a Wails call
 * (Wails rejects with the Go error's string, already French per app.go). */
export function errorMessage(err: unknown): string {
  if (typeof err === "string") return err;
  if (err instanceof Error) return err.message;
  return "Une erreur inattendue est survenue.";
}
