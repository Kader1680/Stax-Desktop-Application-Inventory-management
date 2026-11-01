// Mirrors internal/domain/models.go JSON shape. Kept hand-written (rather than
// relying solely on Wails' generated models) so the frontend compiles even
// before/without a successful bindings generation step.

export type Unit = "piece" | "kg" | "g" | "l";

export interface Product {
  id: number;
  name: string;
  sku: string;
  barcode?: string;
  categoryId?: number | null;
  categoryName?: string;
  unit: Unit;
  purchasePrice: number; // cents
  sellingPrice: number; // cents
  stockQty: string; // decimal string
  minStockQty: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface ProductInput {
  name: string;
  sku: string;
  barcode: string;
  categoryId: number | null;
  unit: string;
  purchasePrice: number; // plain amount, e.g. 19.99 — backend converts to cents
  sellingPrice: number;
  stockQty: string;
  minStockQty: string;
}

export interface StockMovement {
  id: number;
  productId: number;
  type: "in" | "out" | "adjustment" | "sale" | "purchase";
  qty: string;
  reason?: string;
  referenceId?: number | null;
  referenceType?: string;
  createdAt: string;
}

export const UNIT_LABELS: Record<Unit, string> = {
  piece: "Pièce",
  kg: "Kilogramme",
  g: "Gramme",
  l: "Litre",
};

/** Formats a cents integer as a money string, e.g. 1999 -> "19.99". */
export function formatMoney(cents: number): string {
  return (cents / 100).toFixed(2);
}
