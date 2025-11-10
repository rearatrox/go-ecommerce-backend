CREATE TABLE IF NOT EXISTS products (
  id BIGSERIAL PRIMARY KEY,
  sku TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  price_cents INT NOT NULL,
  currency TEXT NOT NULL DEFAULT 'EUR',
  stock_qty INT NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'active',
  image_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
