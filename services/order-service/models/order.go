package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

type Order struct {
	ID                int64       `db:"id" json:"id" swaggerignore:"true"`
	UserID            int64       `db:"user_id" json:"userId" swaggerignore:"true"`
	CartID            int64       `db:"cart_id" json:"cartId" swaggerignore:"true"`
	Status            string      `db:"status" json:"status" example:"pending"`
	TotalCents        int         `db:"total_cents" json:"totalCents" example:"5999"`
	ShippingAddressID *int64      `db:"shipping_address_id" json:"shippingAddressId,omitempty" example:"1"`
	BillingAddressID  *int64      `db:"billing_address_id" json:"billingAddressId,omitempty" example:"1"`
	CreatedAt         time.Time   `db:"created_at" json:"createdAt" swaggerignore:"true"`
	UpdatedAt         *time.Time  `db:"updated_at" json:"updatedAt,omitempty" swaggerignore:"true"`
	Items             []OrderItem `json:"items,omitempty"`

	// Address details (joined)
	ShippingAddress *Address `json:"shippingAddress,omitempty"`
	BillingAddress  *Address `json:"billingAddress,omitempty"`
}

type Address struct {
	ID        int64  `json:"id"`
	Street    string `json:"street" example:"Main Street 123"`
	City      string `json:"city" example:"New York"`
	State     string `json:"state" example:"NY"`
	ZipCode   string `json:"zipCode" example:"10001"`
	Country   string `json:"country" example:"USA"`
	IsDefault bool   `json:"isDefault"`
}

// CreateFromCart creates a new order from an active cart and marks cart as ordered
// used in: handlers.CreateOrder
func CreateFromCart(userId int64, shippingAddressId, billingAddressId *int64) (*Order, error) {
	// Start transaction
	tx, err := db.DB.Begin(db.Ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(db.Ctx)

	// Get active cart
	var cartID int64
	var total int
	err = tx.QueryRow(db.Ctx, `
		SELECT c.id, COALESCE(SUM(ci.price_cents * ci.quantity), 0) as total
		FROM carts c
		LEFT JOIN cart_items ci ON c.id = ci.cart_id
		WHERE c.user_id=$1 AND c.status='active'
		GROUP BY c.id
	`, userId).Scan(&cartID, &total)
	if err != nil {
		return nil, err
	}

	// Create order
	order := &Order{
		UserID:            userId,
		CartID:            cartID,
		Status:            "pending",
		TotalCents:        total,
		ShippingAddressID: shippingAddressId,
		BillingAddressID:  billingAddressId,
	}

	query := `INSERT INTO orders (user_id, cart_id, status, total_cents, shipping_address_id, billing_address_id, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, now())
	          RETURNING id, created_at`
	err = tx.QueryRow(db.Ctx, query, order.UserID, order.CartID, order.Status, order.TotalCents,
		order.ShippingAddressID, order.BillingAddressID).Scan(&order.ID, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Copy cart items to order items
	_, err = tx.Exec(db.Ctx, `
		INSERT INTO order_items (order_id, product_id, quantity, price_cents, product_name, created_at)
		SELECT $1, ci.product_id, ci.quantity, ci.price_cents, p.name, now()
		FROM cart_items ci
		JOIN products p ON ci.product_id = p.id
		WHERE ci.cart_id=$2
	`, order.ID, cartID)
	if err != nil {
		return nil, err
	}

	// Update cart status to 'ordered'
	_, err = tx.Exec(db.Ctx, `UPDATE carts SET status='ordered', updated_at=now() WHERE id=$1`, cartID)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(db.Ctx); err != nil {
		return nil, err
	}

	// Load order items
	if err = order.LoadItems(); err != nil {
		return nil, err
	}

	// Load addresses
	if err = order.LoadAddresses(); err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderByID retrieves a specific order by ID for a user including items and addresses
// used in: handlers.GetOrder, handlers.UpdateOrderStatus
func GetOrderByID(orderId, userId int64) (*Order, error) {
	order := &Order{}
	query := `SELECT id, user_id, cart_id, status, total_cents, shipping_address_id, billing_address_id, created_at, updated_at
	          FROM orders
	          WHERE id=$1 AND user_id=$2`
	err := db.DB.QueryRow(db.Ctx, query, orderId, userId).Scan(
		&order.ID, &order.UserID, &order.CartID, &order.Status, &order.TotalCents,
		&order.ShippingAddressID, &order.BillingAddressID, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Load items and addresses
	if err = order.LoadItems(); err != nil {
		return nil, err
	}
	if err = order.LoadAddresses(); err != nil {
		return nil, err
	}

	return order, nil
}

// GetUserOrders retrieves all orders for a user ordered by creation date
// used in: handlers.ListOrders
func GetUserOrders(userId int64) ([]Order, error) {
	query := `SELECT id, user_id, cart_id, status, total_cents, shipping_address_id, billing_address_id, created_at, updated_at
	          FROM orders
	          WHERE user_id=$1
	          ORDER BY created_at DESC`

	rows, err := db.DB.Query(db.Ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.CartID, &order.Status, &order.TotalCents,
			&order.ShippingAddressID, &order.BillingAddressID, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Load items and addresses for each order
		if err = order.LoadItems(); err != nil {
			return nil, err
		}
		if err = order.LoadAddresses(); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if orders == nil {
		orders = []Order{}
	}

	return orders, nil
}

// LoadItems loads all order items for an order
// used in: CreateFromCart, GetOrderByID, GetUserOrders
func (o *Order) LoadItems() error {
	items, err := GetOrderItems(o.ID)
	if err != nil {
		return err
	}
	o.Items = items
	return nil
}

// LoadAddresses loads shipping and billing address details for an order
// used in: CreateFromCart, GetOrderByID, GetUserOrders
func (o *Order) LoadAddresses() error {
	if o.ShippingAddressID != nil {
		addr := &Address{}
		query := `SELECT id, street, city, state, zip_code, country, is_default
		          FROM addresses WHERE id=$1`
		err := db.DB.QueryRow(db.Ctx, query, *o.ShippingAddressID).Scan(
			&addr.ID, &addr.Street, &addr.City, &addr.State, &addr.ZipCode, &addr.Country, &addr.IsDefault,
		)
		if err == nil {
			o.ShippingAddress = addr
		}
	}

	if o.BillingAddressID != nil {
		addr := &Address{}
		query := `SELECT id, street, city, state, zip_code, country, is_default
		          FROM addresses WHERE id=$1`
		err := db.DB.QueryRow(db.Ctx, query, *o.BillingAddressID).Scan(
			&addr.ID, &addr.Street, &addr.City, &addr.State, &addr.ZipCode, &addr.Country, &addr.IsDefault,
		)
		if err == nil {
			o.BillingAddress = addr
		}
	}

	return nil
}

// UpdateStatus changes the order status (e.g., pending, confirmed, shipped, delivered, cancelled)
// used in: handlers.UpdateOrderStatus
func (o *Order) UpdateStatus(newStatus string) error {
	query := `UPDATE orders SET status=$1, updated_at=now() WHERE id=$2`
	_, err := db.DB.Exec(db.Ctx, query, newStatus, o.ID)
	if err != nil {
		return err
	}
	o.Status = newStatus
	now := time.Now()
	o.UpdatedAt = &now
	return nil
}
