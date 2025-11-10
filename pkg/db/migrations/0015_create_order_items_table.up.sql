-- Create order_items table (snapshot of cart items at order time)
CREATE TABLE IF NOT EXISTS order_items (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  quantity INTEGER NOT NULL CHECK (quantity > 0),
  price_cents INTEGER NOT NULL,
  product_name VARCHAR(255) NOT NULL, -- Snapshot of product name at order time
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);

-- Index for order queries
CREATE INDEX idx_order_items_order_id ON order_items(order_id);

-- Index for product reference
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
