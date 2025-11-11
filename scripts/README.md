# Demo Data Seed Script

This script populates your E-Commerce Backend with demo data for testing and development purposes.

## 📦 What Gets Created

### Users
- **Admin User**
  - Email: `admin@example.com`
  - Password: `Admin123!`
  - Role: admin

- **Customer 1**
  - Email: `customer1@example.com`
  - Password: `Customer123!`
  - Role: customer
  - 2 Addresses (shipping + billing)
  - 3 Items in cart

- **Customer 2**
  - Email: `customer2@example.com`
  - Password: `Customer123!`
  - Role: customer
  - 1 Address
  - 2 Items in cart

### Categories
- Electronics
- Clothing
- Books
- Home & Garden

### Products (10 items)
1. **Laptop Pro 15** - €1,299.00 (Electronics)
2. **Wireless Mouse** - €29.90 (Electronics)
3. **USB-C Cable** - €14.90 (Electronics)
4. **T-Shirt Classic Blue** - €19.99 (Clothing)
5. **Jeans Slim Fit** - €49.99 (Clothing)
6. **Go Programming Language** - €39.99 (Books)
7. **Clean Code** - €44.99 (Books)
8. **LED Desk Lamp** - €34.99 (Home & Garden)
9. **Plant Pot Ceramic** - €12.99 (Home & Garden)
10. **Coffee Mug Set** - €24.99 (Home & Garden)

### Shopping Carts
- Customer 1: Laptop, 2x Mouse, Go Book
- Customer 2: 2x T-Shirt, Coffee Mug Set

## 🚀 Usage

### Prerequisites
- Docker containers must be running: `docker compose up -d`
- All services must be healthy and accessible

### Run the Script

```bash
# From project root
cd scripts
go run seed-demo-data.go
```

Or build and run:

```bash
# Build
go build -o seed-demo-data seed-demo-data.go

# Run
./seed-demo-data
```

### Environment Variables

The script respects your `.env` configuration:

```env
USERSERVICE_PORT=8081
PRODUCTSERVICE_PORT=8082
CARTSERVICE_PORT=8083
API_PREFIX=/api/v1
```

## ⚙️ How It Works

The script executes the following steps:

1. **Wait for Services** - Ensures all services are ready
2. **Create Demo Users** - Registers customer accounts
3. **Login Users** - Obtains JWT tokens for API calls
4. **Create Addresses** - Adds shipping/billing addresses
5. **Create Categories** - Sets up product categories
6. **Create Products** - Adds demo products with stock
7. **Add Cart Items** - Populates shopping carts

## 🔄 Re-running the Script

The script is **idempotent** - you can run it multiple times safely:
- Existing users won't be duplicated (409 Conflict handled)
- Existing categories/products are skipped
- Cart items may be added multiple times (use with caution)

## 🧹 Clean Up

To reset the database and start fresh:

```bash
# Stop containers
docker compose down

# Remove volumes (deletes all data!)
docker volume rm go-ecommerce-backend_db-data

# Start fresh
docker compose up -d

# Wait for migrations, then seed
cd scripts && go run seed-demo-data.go
```

## 📝 Testing Workflow

After seeding, you can test the complete flow:

1. **Login as Customer**
   ```bash
   POST http://localhost:8081/api/v1/users/login
   { "email": "customer1@example.com", "password": "Customer123!" }
   ```

2. **View Cart** (already has items)
   ```bash
   GET http://localhost:8083/api/v1/cart
   Authorization: Bearer <token>
   ```

3. **Create Order**
   ```bash
   POST http://localhost:8084/api/v1/orders
   Authorization: Bearer <token>
   { "shippingAddressId": 1, "billingAddressId": 1 }
   ```

4. **Create Payment Intent**
   ```bash
   POST http://localhost:8085/api/v1/payment-intents
   Authorization: Bearer <token>
   { "orderId": 1 }
   ```

5. **Test with Stripe CLI**
   ```bash
   stripe payment_intents confirm <payment_intent_id> \
     --payment-method pm_card_visa \
     --return-url https://example.com
   ```

## 🐛 Troubleshooting

### "connection refused"
- Ensure `docker compose up -d` is running
- Wait 10-15 seconds for services to fully start
- Check: `docker compose ps`

### "401 Unauthorized" 
- Admin user might not exist yet
- Check logs: `docker compose logs user-service`
- Admin is created by migration `0004_seed_admin_user.up.sql`

### "health endpoint not found"
- Services don't have `/health` endpoints (expected)
- Script will timeout and continue anyway
- Check individual service logs if issues persist

## 💡 Tips

- Use **Swagger UI** to explore created data:
  - Products: http://localhost:8082/api/v1/products/swagger/index.html
  - Users: http://localhost:8081/api/v1/users/swagger/index.html
  - Cart: http://localhost:8083/api/v1/cart/swagger/index.html
  - Orders: http://localhost:8084/api/v1/orders/swagger/index.html
  - Payments: http://localhost:8085/api/v1/payments/swagger/index.html

- Use **pgweb** to inspect the database:
  - URL: http://localhost:8088

## 🔐 Security Note

**⚠️ This script is for development/demo purposes only!**

- Passwords are hardcoded and simple
- No rate limiting or security checks
- Do not use in production!
