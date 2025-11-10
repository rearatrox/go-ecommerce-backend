-- Seed initial admin user
INSERT INTO users (email, password, role)
VALUES (
  'admin@example.com',
  '$2a$10$MC9YbBqsBZ2HLM1Lzj8aReVzIFHP.44PoZ1wQ6e/nEpGc4zWM/8qG', -- bcrypt("admin123")
  'admin'
)
ON CONFLICT (email) DO NOTHING;
