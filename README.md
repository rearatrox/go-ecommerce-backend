# Go E-Commerce Backend

A modular **E-Commerce Backend in Go**, consisting of multiple microservices.  
Currently, the project includes the following services:

- **User-Service** – Authentication, registration, profile management, and addresses  
- **Product-Service** – Product management with categories and many-to-many relationships  
- **Cart-Service** – Shopping cart management with automatic price snapshot functionality  
- **Order-Service** – Order processing and history with address snapshots
- **Payment-Service** – Stripe payment integration with webhooks and retry logic

Each service runs as an independent container in the Docker Compose setup and uses a shared PostgreSQL database with automatic migrations.

---

## Features

### 🏗️ Architecture
- Clean **Microservices Architecture** in Go with `gin-gonic`
- Shared `.env` configuration (via `.env.example`)
- Multi-service setup with **Docker Compose**
- **Automatic database migrations** with golang-migrate
- **Internal service authentication** with shared secrets
- Ready for future **Kubernetes deployments**
- Each service has its own **Swagger documentation**

### 🔐 Security
- **JWT-based authentication** with Role-Based Access Control (RBAC)
- Admin-protected routes with middleware
- Password hashing with bcrypt
- Token version management for secure logout functionality
- Internal API endpoints protected with shared secret authentication
- Address and payment ownership validation

### 📦 Product-Service
- CRUD operations for products (with SKU, prices in cents, stock management)
- Category system with slug-based routing
- Many-to-many relationship between products and categories
- Admin-only product management

### 👤 User-Service
- Registration and login with JWT
- Profile management (first name, last name, phone)
- Address management (shipping/billing addresses)
- Automatic default address management
- Admin area for user management

### 🛒 Cart-Service
- Automatic cart creation and management
- One active cart per user (via UNIQUE constraint)
- Price snapshot when adding items (protects against price changes)
- Automatic quantity merging when adding duplicates
- Status management (active, ordered, abandoned)
- Join with product data for complete item information

### 📦 Order-Service
- Create orders from active cart with automatic status management
- Order history with complete item and address details
- Price and product name snapshots at order time
- Status tracking (pending, confirmed, shipped, delivered, cancelled)
- Address linking (shipping and billing)
- Address ownership validation for security
- Order cancellation with stock restoration
- Automatic stock reduction when orders are confirmed

### 💳 Payment-Service
- **Stripe integration** with Payment Intents API
- Secure webhook handling with signature verification
- Payment retry logic for failed/cancelled payments
- Order ownership validation before payment creation
- Automatic order status updates after successful payment
- Status management (pending, processing, succeeded, failed, cancelled, superseded)
- Webhook-triggered stock reduction on successful payments

### �🛠️ Developer Experience
- Structured **logging** with slog and context propagation
- **Hot-reload** possible during development
- Integrated **pgweb** for database inspection (port 8088)
- Swagger UI for all services with live testing

---

## Installation & Setup

1. **Clone repository**
   ```
   git clone https://github.com/rearatrox/go-ecommerce-backend.git
   cd go-ecommerce-backend
   ```

2. **Adjust .env files**  
   Create a `.env` file from `.env.example` and customize it:
   ```
   cp .env.example .env
   ```

3. **Start containers**
   ```
   docker compose up -d
   ```

4. **Test services**   
   - User-Service: [http://localhost:8081](http://localhost:8081)
   - Product-Service: [http://localhost:8082](http://localhost:8082)
   - Cart-Service: [http://localhost:8083](http://localhost:8083)
   - Order-Service: [http://localhost:8084](http://localhost:8084)
   - Payment-Service: [http://localhost:8085](http://localhost:8085)
   - pgweb (DB-Admin): [http://localhost:8088](http://localhost:8088)

5. **Open Swagger UI**
   - User-Service Swagger: [http://localhost:8081/api/v1/users/swagger/index.html](http://localhost:8081/api/v1/users/swagger/index.html)
   - Product-Service Swagger: [http://localhost:8082/api/v1/products/swagger/index.html](http://localhost:8082/api/v1/products/swagger/index.html)
   - Cart-Service Swagger: [http://localhost:8083/api/v1/cart/swagger/index.html](http://localhost:8083/api/v1/cart/swagger/index.html)
   - Order-Service Swagger: [http://localhost:8084/api/v1/orders/swagger/index.html](http://localhost:8084/api/v1/orders/swagger/index.html)
   - Payment-Service Swagger: [http://localhost:8085/api/v1/payments/swagger/index.html](http://localhost:8085/api/v1/payments/swagger/index.html) 

---

## ⚙️ Environment Variables (`.env.example`)

| Variable | Description | Example Value |
|-----------|---------------|---------------|
| **API_PREFIX** | Common API prefix for all services | `/api/v1` |
| **JWT_SECRET** | Secret key for JWT token signing | `supersecret` |
| **INTERNAL_API_SECRET** | Shared secret for internal service-to-service communication | `internal-secret-key` |

### 💳 Stripe

| Variable | Description | Example Value |
|-----------|---------------|---------------|
| **STRIPE_SECRET_KEY** | Stripe API secret key | `sk_test_...` |
| **STRIPE_WEBHOOK_SECRET** | Stripe webhook signing secret | `whsec_...` |

### 🪵 Logger

| Variable | Description | Example Value |
|-----------|---------------|---------------|
| **LOG_LEVEL** | Log level (e.g. `debug`, `info`, `warn`, `error`) | `info` |
| **LOG_FORMAT** | Log format (`text` or `json`) | `json` |
| **LOG_OUTPUT** | Log output destination (`stdout`, `file`, etc.) | `stdout` |
| **REQUEST_ID_HEADER** | Header name for request IDs (tracing) | `X-Request-Id` |

### 🧩 Services

| Variable | Description | Example Value |
|-----------|---------------|---------------|
| **USERSERVICE_PORT** | External port of User-Service | `8081` |
| **PRODUCTSERVICE_PORT** | External port of Product-Service | `8082` |
| **CARTSERVICE_PORT** | External port of Cart-Service | `8083` |
| **ORDERSERVICE_PORT** | External port of Order-Service | `8084` |
| **PAYMENTSERVICE_PORT** | External port of Payment-Service | `8085` |

### 🗄️ Database

| Variable | Description | Example Value |
|-----------|---------------|---------------|
| **DB_HOST** | Hostname (must match Docker Service-Name!) for PostgreSQL | `api-database` |
| **DB_USERNAME** | Username for PostgreSQL | `admin` |
| **DB_PASSWORD** | Password for PostgreSQL | `password123` |
| **DB_NAME** | Database name | `api_db` |
| **DB_PORT** | Port of PostgreSQL instance | `5432` |
| **DB_SSLMODE** | SSL mode of connection (`disable`, `require`, etc.) | `disable` |

> 💡 **Note:**  
> The DATABASE_URL is automatically generated with the above settings

---

## 📘 Swagger API Documentation

Each service has its own Swagger documentation based on [swaggo/gin-swagger](https://github.com/swaggo/gin-swagger).

The Swagger files are automatically generated during build and enable interactive documentation of all API endpoints.

### 👤 User-Service

- **Port:** `${USERSERVICE_PORT}` (default: `8081`)  
- **Swagger-URL:** [http://localhost:8081/api/v1/users/swagger/index.html](http://localhost:8081/api/v1/users/swagger/index.html)

### 📦 Product-Service

- **Port:** `${PRODUCTSERVICE_PORT}` (default: `8082`)  
- **Swagger-URL:** [http://localhost:8082/api/v1/products/swagger/index.html](http://localhost:8082/api/v1/products/swagger/index.html)

### 🛒 Cart-Service

- **Port:** `${CARTSERVICE_PORT}` (default: `8083`)  
- **Swagger-URL:** [http://localhost:8083/api/v1/cart/swagger/index.html](http://localhost:8083/api/v1/cart/swagger/index.html)

### 📦 Order-Service

- **Port:** `${ORDERSERVICE_PORT}` (default: `8084`)  
- **Swagger-URL:** [http://localhost:8084/api/v1/orders/swagger/index.html](http://localhost:8084/api/v1/orders/swagger/index.html)

### 💳 Payment-Service

- **Port:** `${PAYMENTSERVICE_PORT}` (default: `8085`)  
- **Swagger-URL:** [http://localhost:8085/api/v1/payments/swagger/index.html](http://localhost:8085/api/v1/payments/swagger/index.html)


> 💡 **Authentication:**  
> Protected endpoints require a JWT token in the `Authorization` header: `Bearer <token>`  
> You receive the token after successful login via `/api/v1/auth/login`

> 💡 **Note:**  
> Ports are dynamically set via the respective ENV variables,  
> so the Swagger UI automatically uses the correct host in any environment (local or container).

---

## 🗄️ Database Structure

The project uses **PostgreSQL** with automatic migrations via `golang-migrate`.

### Tables

**User-Service:**
- `users` - Users with email, password (bcrypt), role, token version, personal info
- `addresses` - Shipping and billing addresses with default management

**Product-Service:**
- `products` - Products with SKU, name, price (in cents), stock, status, images
- `categories` - Categories with slug for SEO-friendly URLs
- `product_categories` - Junction table for many-to-many relationship

**Cart-Service:**
- `carts` - Shopping carts with user assignment and status (active/ordered/abandoned)
- `cart_items` - Products in cart with quantity and price snapshot

**Order-Service:**
**Order-Service:**
- `orders` - Orders with status, total, and address references
- `order_items` - Order items with product snapshots (name, price) at order time

**Payment-Service:**
- `payments` - Payment records with Stripe integration, status tracking, and order linkage

### Migrations

All migrations are consolidated in a single initial schema file under `/pkg/db/migrations/`:

```bash
0001_initial_schema.up.sql     # Complete database schema with all tables
0001_initial_schema.down.sql   # Rollback for complete schema
```

The consolidated migration includes:
- All tables (users, addresses, products, categories, carts, orders, payments)
- All relationships and foreign keys
- All indexes for performance optimization
- Default admin user (email: `admin@example.com`, password: `admin123`)
- All constraints and data types

> 💡 **pgweb:**  
> You can visually inspect the database via pgweb: [http://localhost:8088](http://localhost:8088)

---

## 🎲 Demo Data

For testing purposes, you can use the included seed script to populate the database with demo data:

```bash
cd scripts
go run seed-demo-data.go
```

The script creates:
- **3 demo users** (admin, customer1, customer2)
- **2 addresses per user** (shipping & billing)
- **4 product categories** (Electronics, Clothing, Home & Garden, Sports)
- **10 sample products** with stock and images
- **Pre-filled shopping carts** for each customer

**Demo Users:**
- `admin@example.com` - Admin user (password: `admin123`)
- `customer1@demo.com` - Customer with shopping cart
- `customer2@demo.com` - Customer with shopping cart

> 💡 **Note:** Make sure all services are running before executing the seed script!

For more information, see `/scripts/README.md`

---

## 🚀 Roadmap

### ✅ Implemented
- [x] User-Service with auth and profile management
- [x] Product-Service with categories
- [x] Cart-Service with price snapshots
- [x] Order-Service with address snapshots and status management
- [x] Payment-Service with Stripe integration
- [x] JWT-based authentication
- [x] Role-Based Access Control (Admin/User)
- [x] Automatic database migrations
- [x] Swagger documentation for all services
- [x] Stock validation - Check if enough inventory available when adding to cart
- [x] Inventory management - Stock reduction for successful orders
- [x] Payment retry logic - Retry failed/cancelled payments
- [x] Webhook integration - Stripe webhook handling with signature verification
- [x] Internal API security - Service-to-service authentication
- [x] Order cancellation - Cancel orders with stock restoration

### 🔄 Planned (Priority)
- [ ] Search and filter functions - Filter products by criteria
- [ ] Pagination - For large product lists
- [ ] PayPal integration - Additional payment provider

### 💡 Nice-to-Have
- [ ] Review/rating system for products
- [ ] Wishlist functionality
- [ ] Email verification and password reset
- [ ] Refresh tokens
- [ ] Notification service
- [ ] Admin dashboard with analytics
- [ ] API gateway (Kong/Traefik)
- [ ] Image upload (S3/MinIO)

---

## ☸️ Kubernetes 

In the future, Kubernetes manifests will be provided under  
`/k8s/` to enable easy deployment of services on a cluster.

---

## 📚 Foundation & License

This project was based on the following Udemy course: [Go - The Complete Guide](https://www.udemy.com/course/go-the-complete-guide/)

**License:** MIT (see LICENSE file)

---

## 🤝 Contributing

Contributions are welcome! Please create a pull request or open an issue for suggestions.

**Project structure:**
```
go-ecommerce-backend/
├── pkg/                          # Shared packages
│   ├── db/                       # Database connection & migrations
│   ├── logger/                   # Structured logging
│   └── middleware/
│       ├── auth/                 # JWT auth middleware
│       └── serviceauth/          # Internal service authentication
├── services/
│   ├── user-service/             # User, Auth, Addresses
│   ├── product-service/          # Products, Categories
│   ├── cart-service/             # Shopping cart
│   ├── order-service/            # Orders & order history
│   └── payment-service/          # Stripe payment integration
├── scripts/                      # Utility scripts
│   └── seed-demo-data.go        # Demo data seeding tool
├── docker-compose.yaml           # Multi-service setup
├── .env.example                  # Environment template
└── README.md
```

---

## 📧 Contact

**Author:** Tim Hauschild  
**Website:** [webdesign-hauschild.de](https://webdesign-hauschild.de)  
**GitHub:** [rearatrox](https://github.com/rearatrox)
