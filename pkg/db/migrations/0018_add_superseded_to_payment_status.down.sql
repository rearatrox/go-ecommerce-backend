-- Remove 'superseded' from payment status enum
ALTER TABLE payments DROP CONSTRAINT IF EXISTS payments_status_check;
ALTER TABLE payments ADD CONSTRAINT payments_status_check 
  CHECK (status IN ('pending', 'processing', 'succeeded', 'failed', 'cancelled'));
