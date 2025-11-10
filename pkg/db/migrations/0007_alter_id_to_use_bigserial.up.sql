-- Alter existing integer id column to use BIGSERIAL (64-bit auto-increment)

-- 1️⃣  Entferne alte Default-Sequenz, falls vorhanden
ALTER TABLE products ALTER COLUMN id DROP DEFAULT;

-- 2️⃣  Setze Datentyp auf BIGINT (64-bit)
ALTER TABLE products ALTER COLUMN id TYPE BIGINT;

-- 3️⃣  Erstelle eine neue Sequenz falls keine existiert
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_class WHERE relname = 'products_id_seq'
  ) THEN
    CREATE SEQUENCE products_id_seq OWNED BY products.id;
  END IF;
END $$;

-- 4️⃣  Hänge die neue Sequenz als Default dran
ALTER TABLE products ALTER COLUMN id SET DEFAULT nextval('products_id_seq');

-- 5️⃣  Synchronisiere die Sequenz mit dem aktuellen Maximum
SELECT setval('products_id_seq', COALESCE((SELECT MAX(id) FROM products), 1), true);
