-- Create carts table (one active cart per user)
CREATE TABLE IF NOT EXISTS carts (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'ordered', 'abandoned')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ
);

-- Only one active cart per user
CREATE UNIQUE INDEX idx_carts_user_active ON carts(user_id) WHERE status = 'active';

-- Index for status queries
CREATE INDEX idx_carts_status ON carts(status);
