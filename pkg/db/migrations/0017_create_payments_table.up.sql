-- Create payments table
CREATE TABLE IF NOT EXISTS payments (
  id BIGSERIAL PRIMARY KEY,
  order_id BIGINT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount_cents INTEGER NOT NULL CHECK (amount_cents > 0),
  currency VARCHAR(3) NOT NULL DEFAULT 'EUR',
  status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'succeeded', 'failed', 'cancelled')),
  stripe_payment_intent_id VARCHAR(255) UNIQUE,
  stripe_client_secret VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);

-- Index for order payments
CREATE INDEX idx_payments_order_id ON payments(order_id);

-- Index for user payments
CREATE INDEX idx_payments_user_id ON payments(user_id);

-- Index for status queries
CREATE INDEX idx_payments_status ON payments(status);

-- Index for Stripe payment intent lookups
CREATE INDEX idx_payments_stripe_payment_intent_id ON payments(stripe_payment_intent_id);
