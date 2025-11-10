-- Create orders table
CREATE TABLE IF NOT EXISTS orders (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  cart_id BIGINT NOT NULL REFERENCES carts(id) ON DELETE CASCADE,
  status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'shipped', 'delivered', 'cancelled')),
  total_cents INTEGER NOT NULL,
  shipping_address_id BIGINT REFERENCES addresses(id) ON DELETE SET NULL,
  billing_address_id BIGINT REFERENCES addresses(id) ON DELETE SET NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);

-- Index for user orders
CREATE INDEX idx_orders_user_id ON orders(user_id);

-- Index for status queries
CREATE INDEX idx_orders_status ON orders(status);

-- Index for cart reference
CREATE INDEX idx_orders_cart_id ON orders(cart_id);
