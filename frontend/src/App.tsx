import { HashRouter, Routes, Route, NavLink } from "react-router-dom";
import fr from "./i18n/fr.json";
import ProductsPage from "./pages/ProductsPage";
import DashboardPage from "./pages/DashboardPage";
import InventoryPage from "./pages/InventoryPage";
import WelcomePage from "./pages/WelcomePage";

function Placeholder({ title }: { title: string }) {
  return (
    <div className="card">
      <h2>{title}</h2>
      <p>Module à venir.</p>
    </div>
  );
}

export default function App() {
  return (
    <HashRouter>
      <div className="app-shell">
        <aside className="sidebar">
          <nav>
            <NavLink to="/welcome" end>weclome</NavLink>
            <NavLink to="/" end>{fr.nav.dashboard}</NavLink>
            <NavLink to="/products">{fr.nav.products}</NavLink>
            <NavLink to="/inventory">{fr.nav.inventory}</NavLink>
            <NavLink to="/pos">{fr.nav.pos}</NavLink>
            <NavLink to="/purchases">{fr.nav.purchases}</NavLink>
            <NavLink to="/customers">{fr.nav.customers}</NavLink>
            <NavLink to="/suppliers">{fr.nav.suppliers}</NavLink>
            <NavLink to="/reports">{fr.nav.reports}</NavLink>
          </nav>
        </aside>
        <main className="content">
          <Routes>
            <Route path="/" element={<DashboardPage />} />
            <Route path="/welcome" element={< WelcomePage/>} />
            <Route path="/products" element={<ProductsPage />} />
            <Route path="/inventory" element={<InventoryPage />} />
            <Route path="/pos" element={<Placeholder title={fr.nav.pos} />} />
            <Route path="/purchases" element={<Placeholder title={fr.nav.purchases} />} />
            <Route path="/customers" element={<Placeholder title={fr.nav.customers} />} />
            <Route path="/suppliers" element={<Placeholder title={fr.nav.suppliers} />} />
            <Route path="/reports" element={<Placeholder title={fr.nav.reports} />} />
          </Routes>
        </main>
      </div>
    </HashRouter>
  );
}
