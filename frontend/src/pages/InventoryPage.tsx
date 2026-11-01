import { useEffect, useState } from "react";
import { api, errorMessage } from "../lib/api";
import { Product, StockMovement } from "../lib/types";

const MOVEMENT_LABELS: Record<string, string> = {
  in: "Entrée",
  out: "Sortie",
  adjustment: "Ajustement",
  sale: "Vente",
  purchase: "Achat",
};

export default function InventoryPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [productId, setProductId] = useState<number | null>(null);
  const [mode, setMode] = useState<"in" | "out" | "adjustment">("in");
  const [qty, setQty] = useState("");
  const [reason, setReason] = useState("");
  const [movements, setMovements] = useState<StockMovement[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);

  useEffect(() => {
    api.searchProducts("", null).then((list) => setProducts(list ?? [])).catch((e) => setError(errorMessage(e)));
    refreshMovements();
  }, []);

  async function refreshMovements() {
    try {
      setMovements((await api.recentStockMovements(50)) ?? []);
    } catch (e) {
      setError(errorMessage(e));
    }
  }

  async function submit() {
    setError(null);
    setMessage(null);
    if (!productId) {
      setError("Veuillez sélectionner un produit.");
      return;
    }
    try {
      if (mode === "in") await api.stockIn(productId, qty, reason);
      else if (mode === "out") await api.stockOut(productId, qty, reason);
      else await api.stockAdjustment(productId, qty, reason);
      setMessage("Mouvement de stock enregistré.");
      setQty("");
      setReason("");
      await refreshMovements();
      const refreshed = await api.searchProducts("", null);
      setProducts(refreshed ?? []);
    } catch (e) {
      setError(errorMessage(e));
    }
  }

  return (
    <div>
      <h2>Inventaire</h2>

      <div className="card" style={{ marginBottom: 16 }}>
        <div style={{ display: "flex", gap: 12, flexWrap: "wrap", alignItems: "end" }}>
          <label style={{ display: "flex", flexDirection: "column", fontSize: 13 }}>
            Produit
            <select value={productId ?? ""} onChange={(e) => setProductId(Number(e.target.value) || null)}>
              <option value="">-- Sélectionner --</option>
              {products.map((p) => (
                <option key={p.id} value={p.id}>{p.name} ({p.sku})</option>
              ))}
            </select>
          </label>
          <label style={{ display: "flex", flexDirection: "column", fontSize: 13 }}>
            Type
            <select value={mode} onChange={(e) => setMode(e.target.value as any)}>
              <option value="in">Entrée de stock</option>
              <option value="out">Sortie de stock</option>
              <option value="adjustment">Ajustement (correction)</option>
            </select>
          </label>
          <label style={{ display: "flex", flexDirection: "column", fontSize: 13 }}>
            {mode === "adjustment" ? "Delta (+/-)" : "Quantité"}
            <input value={qty} onChange={(e) => setQty(e.target.value)} placeholder="ex: 10 ou 0.5" />
          </label>
          <label style={{ display: "flex", flexDirection: "column", fontSize: 13, flex: 1, minWidth: 200 }}>
            Raison
            <input value={reason} onChange={(e) => setReason(e.target.value)} placeholder="ex: inventaire physique" />
          </label>
          <button className="btn" onClick={submit}>Enregistrer</button>
        </div>
        {error && <p className="error-text">{error}</p>}
        {message && <p style={{ color: "green", fontSize: 13 }}>{message}</p>}
      </div>

      <div className="card">
        <h3>Mouvements récents</h3>
        <table>
          <thead>
            <tr>
              <th>Date</th>
              <th>Produit</th>
              <th>Type</th>
              <th>Qté</th>
              <th>Raison</th>
            </tr>
          </thead>
          <tbody>
            {movements.map((m) => (
              <tr key={m.id}>
                <td>{new Date(m.createdAt).toLocaleString("fr-FR")}</td>
                <td>{products.find((p) => p.id === m.productId)?.name ?? m.productId}</td>
                <td>{MOVEMENT_LABELS[m.type] ?? m.type}</td>
                <td>{m.qty}</td>
                <td>{m.reason}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
