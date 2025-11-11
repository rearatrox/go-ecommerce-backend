package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

type OrderItem struct {
	ID          int64      `db:"id" json:"id" swaggerignore:"true"`
	OrderID     int64      `db:"order_id" json:"orderId" swaggerignore:"true"`
	ProductID   int64      `db:"product_id" json:"productId" example:"1"`
	Quantity    int        `db:"quantity" json:"quantity" example:"2"`
	PriceCents  int        `db:"price_cents" json:"priceCents" example:"2999"`
	ProductName string     `db:"product_name" json:"productName" example:"Gaming Laptop"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt" swaggerignore:"true"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty" swaggerignore:"true"`
}

// GetOrderItems retrieves all items for a specific order with product details
// used in: order.LoadItems
func GetOrderItems(orderId int64) ([]OrderItem, error) {
	query := `SELECT id, order_id, product_id, quantity, price_cents, product_name, created_at, updated_at
	          FROM order_items
	          WHERE order_id=$1
	          ORDER BY created_at DESC`

	rows, err := db.DB.Query(db.Ctx, query, orderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.PriceCents,
			&item.ProductName, &item.CreatedAt, &item.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if items == nil {
		items = []OrderItem{}
	}

	return items, nil
}

type CartItem struct {
	ProductID   int64  `db:"product_id"`
	Quantity    int    `db:"quantity"`
	ProductName string `db:"product_name"`
}

// GetCartItemsForUser retrieves cart items for stock validation before creating an order
// used in: handlers.CreateOrder
func GetCartItemsForUser(userId int64) ([]CartItem, error) {
	query := `SELECT ci.product_id, ci.quantity, p.name
	          FROM cart_items ci
	          JOIN carts c ON ci.cart_id = c.id
	          JOIN products p ON ci.product_id = p.id
	          WHERE c.user_id=$1 AND c.status='active'`

	rows, err := db.DB.Query(db.Ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var item CartItem
		err := rows.Scan(&item.ProductID, &item.Quantity, &item.ProductName)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if items == nil {
		items = []CartItem{}
	}

	return items, nil
}
