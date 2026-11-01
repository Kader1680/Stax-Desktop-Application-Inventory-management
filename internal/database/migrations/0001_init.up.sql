CREATE TABLE categories (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL UNIQUE,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE products (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT NOT NULL,
    sku             TEXT NOT NULL UNIQUE,
    barcode         TEXT UNIQUE,
    category_id     INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    unit            TEXT NOT NULL CHECK (unit IN ('piece','kg','g','l')),
    purchase_price  INTEGER NOT NULL DEFAULT 0, -- cents
    selling_price   INTEGER NOT NULL DEFAULT 0, -- cents
    stock_qty       TEXT NOT NULL DEFAULT '0',  -- decimal string
    min_stock_qty   TEXT NOT NULL DEFAULT '0',
    is_active       INTEGER NOT NULL DEFAULT 1,
    is_demo         INTEGER NOT NULL DEFAULT 0,
    created_at      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now')),
    updated_at      TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_barcode ON products(barcode);
CREATE INDEX idx_products_category ON products(category_id);

CREATE TABLE stock_movements (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id     INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    type           TEXT NOT NULL CHECK (type IN ('in','out','adjustment','sale','purchase')),
    qty            TEXT NOT NULL, -- signed decimal string
    reason         TEXT,
    reference_id   INTEGER,
    reference_type TEXT,
    created_at     TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_movements_product ON stock_movements(product_id);
CREATE INDEX idx_movements_created ON stock_movements(created_at);

CREATE TABLE customers (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    phone      TEXT,
    address    TEXT,
    is_demo    INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE suppliers (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    phone      TEXT,
    address    TEXT,
    is_demo    INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE sales (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    invoice_no  TEXT NOT NULL UNIQUE,
    customer_id INTEGER REFERENCES customers(id) ON DELETE SET NULL,
    subtotal    INTEGER NOT NULL,
    discount    INTEGER NOT NULL DEFAULT 0,
    total       INTEGER NOT NULL,
    status      TEXT NOT NULL CHECK (status IN ('completed','cancelled')) DEFAULT 'completed',
    is_demo     INTEGER NOT NULL DEFAULT 0,
    created_at  TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);
CREATE INDEX idx_sales_created ON sales(created_at);

CREATE TABLE sale_items (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    sale_id    INTEGER NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
    product_id INTEGER NOT NULL REFERENCES products(id),
    qty        TEXT NOT NULL,
    unit_price INTEGER NOT NULL,
    discount   INTEGER NOT NULL DEFAULT 0,
    line_total INTEGER NOT NULL
);
CREATE INDEX idx_sale_items_sale ON sale_items(sale_id);
CREATE INDEX idx_sale_items_product ON sale_items(product_id);

CREATE TABLE purchases (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    supplier_id  INTEGER REFERENCES suppliers(id) ON DELETE SET NULL,
    invoice_ref  TEXT,
    total        INTEGER NOT NULL,
    is_demo      INTEGER NOT NULL DEFAULT 0,
    created_at   TEXT NOT NULL DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ','now'))
);

CREATE TABLE purchase_items (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_id   INTEGER NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    product_id    INTEGER NOT NULL REFERENCES products(id),
    qty           TEXT NOT NULL,
    unit_cost     INTEGER NOT NULL,
    line_total    INTEGER NOT NULL
);
CREATE INDEX idx_purchase_items_purchase ON purchase_items(purchase_id);

CREATE TABLE settings (
    key   TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
