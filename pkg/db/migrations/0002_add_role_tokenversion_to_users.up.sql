-- Add user roles and token version tracking
ALTER TABLE users
  ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user',
  ADD COLUMN IF NOT EXISTS token_version INTEGER NOT NULL DEFAULT 0;
