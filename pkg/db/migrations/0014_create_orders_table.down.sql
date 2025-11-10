-- Drop orders table
DROP INDEX IF EXISTS idx_orders_cart_id;
DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_user_id;
DROP TABLE IF EXISTS orders;
