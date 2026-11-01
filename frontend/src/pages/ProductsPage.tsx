import { ReactNode, useEffect, useState } from "react";
import { api, errorMessage } from "../lib/api";
import { Product, ProductInput, UNIT_LABELS, formatMoney, Unit } from "../lib/types";
import fr from "../i18n/fr.json";

const emptyForm: ProductInput = {
  name: "",
  sku: "",
  barcode: "",
  categoryId: null,
  unit: "piece",
  purchasePrice: 0,
  sellingPrice: 0,
  stockQty: "0",
  minStockQty: "0",
};

export default function ProductsPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [query, setQuery] = useState("");
  const [editingId, setEditingId] = useState<number | null>(null);
  const [form, setForm] = useState<ProductInput>(emptyForm);
  const [showForm, setShowForm] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function load(q = query) {
    setLoading(true);
    try {
      const list = await api.searchProducts(q, null);
      setProducts(list ?? []);
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    load("");
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function startCreate() {
    setForm(emptyForm);
    setEditingId(null);
    setShowForm(true);
    setError(null);
  }

  function startEdit(p: Product) {
    setForm({
      name: p.name,
      sku: p.sku,
      barcode: p.barcode ?? "",
      categoryId: p.categoryId ?? null,
      unit: p.unit,
      purchasePrice: p.purchasePrice / 100,
      sellingPrice: p.sellingPrice / 100,
      stockQty: p.stockQty,
      minStockQty: p.minStockQty,
    });
    setEditingId(p.id);
    setShowForm(true);
    setError(null);
  }

  async function save() {
    setError(null);
    try {
      if (editingId) {
        await api.updateProduct(editingId, form);
      } else {
        await api.createProduct(form);
      }
      setShowForm(false);
      await load();
    } catch (err) {
      setError(errorMessage(err));
    }
  }

  async function remove(id: number) {
    if (!confirm("Supprimer ce produit ?")) return;
    try {
      await api.deleteProduct(id);
      await load();
    } catch (err) {
      setError(errorMessage(err));
    }
  }

  return (
    <div>
      <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 16 }}>
        <h2>{fr.products.title}</h2>
        <button className="btn" onClick={startCreate}>{fr.products.new}</button>
      </div>

      <div className="card" style={{ marginBottom: 16 }}>
        <input
          placeholder={fr.products.search}
          value={query}
          onChange={(e) => {
            setQuery(e.target.value);
            load(e.target.value);
          }}
          style={{ width: "100%", padding: 8 }}
        />
      </div>

      {error && <p className="error-text">{error}</p>}

      {showForm && (
        <div className="card" style={{ marginBottom: 16 }}>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 12 }}>
            <Field label={fr.products.name}>
              <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} />
            </Field>
            <Field label={fr.products.sku}>
              <input value={form.sku} onChange={(e) => setForm({ ...form, sku: e.target.value })} />
            </Field>
            <Field label={fr.products.barcode}>
              <input value={form.barcode} onChange={(e) => setForm({ ...form, barcode: e.target.value })} />
            </Field>
            <Field label={fr.products.unit}>
              <select value={form.unit} onChange={(e) => setForm({ ...form, unit: e.target.value as Unit })}>
                {Object.entries(UNIT_LABELS).map(([val, label]) => (
                  <option key={val} value={val}>{label}</option>
                ))}
              </select>
            </Field>
            <Field label={fr.products.purchasePrice}>
              <input
                type="number" step="0.01" value={form.purchasePrice}
                onChange={(e) => setForm({ ...form, purchasePrice: parseFloat(e.target.value) || 0 })}
              />
            </Field>
            <Field label={fr.products.sellingPrice}>
              <input
                type="number" step="0.01" value={form.sellingPrice}
                onChange={(e) => setForm({ ...form, sellingPrice: parseFloat(e.target.value) || 0 })}
              />
            </Field>
            {!editingId && (
              <Field label={fr.products.stock}>
                <input value={form.stockQty} onChange={(e) => setForm({ ...form, stockQty: e.target.value })} />
              </Field>
            )}
            <Field label={fr.products.minStock}>
              <input value={form.minStockQty} onChange={(e) => setForm({ ...form, minStockQty: e.target.value })} />
            </Field>
          </div>
          <div style={{ marginTop: 12 }}>
            <button className="btn" onClick={save}>{fr.products.save}</button>{" "}
            <button onClick={() => setShowForm(false)}>{fr.products.cancel}</button>
          </div>
        </div>
      )}

      <div className="card">
        {loading ? (
          <p>Chargement...</p>
        ) : (
          <table>
            <thead>
              <tr>
                <th>{fr.products.name}</th>
                <th>{fr.products.sku}</th>
                <th>{fr.products.unit}</th>
                <th>{fr.products.sellingPrice}</th>
                <th>{fr.products.stock}</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {products.map((p) => (
                <tr key={p.id}>
                  <td>{p.name}</td>
                  <td>{p.sku}</td>
                  <td>{UNIT_LABELS[p.unit]}</td>
                  <td>{formatMoney(p.sellingPrice)}</td>
                  <td className={Number(p.stockQty) <= Number(p.minStockQty) && Number(p.minStockQty) > 0 ? "low-stock" : ""}>
                    {p.stockQty}
                  </td>
                  <td>
                    <button onClick={() => startEdit(p)}>Modifier</button>{" "}
                    <button onClick={() => remove(p.id)}>{fr.products.delete}</button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <label style={{ display: "flex", flexDirection: "column", fontSize: 13, gap: 4 }}>
      {label}
      {children}
    </label>
  );
}
