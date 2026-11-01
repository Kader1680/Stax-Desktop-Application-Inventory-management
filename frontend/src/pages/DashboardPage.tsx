import { useEffect, useState } from "react";
import { api, errorMessage } from "../lib/api";
import { Product } from "../lib/types";
import fr from "../i18n/fr.json";

export default function DashboardPage() {
  const [lowStock, setLowStock] = useState<Product[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    api.lowStockProducts().then((list) => setLowStock(list ?? [])).catch((e) => setError(errorMessage(e)));
  }, []);

  return (
    <div>
      <h2>{fr.nav.dashboard}</h2>
      <p style={{ color: "#666", fontSize: 13 }}>
        Les ventes du jour, le chiffre d'affaires et le graphique arriveront avec le module Point de vente.
      </p>

      <div className="card">
        <h3>{fr.products.lowStock}</h3>
        {error && <p className="error-text">{error}</p>}
        {lowStock.length === 0 ? (
          <p>Aucun produit en stock faible.</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>{fr.products.name}</th>
                <th>{fr.products.stock}</th>
                <th>{fr.products.minStock}</th>
              </tr>
            </thead>
            <tbody>
              {lowStock.map((p) => (
                <tr key={p.id}>
                  <td>{p.name}</td>
                  <td className="low-stock">{p.stockQty}</td>
                  <td>{p.minStockQty}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}
