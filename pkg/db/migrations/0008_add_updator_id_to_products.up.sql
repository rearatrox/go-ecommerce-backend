-- Add updator_id column to products table
ALTER TABLE products
ADD COLUMN IF NOT EXISTS updator_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
