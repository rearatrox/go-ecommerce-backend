package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

type Cart struct {
	ID        int64      `db:"id" json:"id" swaggerignore:"true"`
	UserID    int64      `db:"user_id" json:"userId" swaggerignore:"true"`
	Status    string     `db:"status" json:"status" example:"active"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt" swaggerignore:"true"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty" swaggerignore:"true"`
	Items     []CartItem `json:"items,omitempty"`
	Total     int        `json:"totalCents,omitempty"` // Total in cents
}

// GetOrCreateCart retrieves the active cart for a user or creates a new one if none exists
// used in: handlers.GetCart, handlers.AddItem, handlers.UpdateItem, handlers.RemoveItem, handlers.ClearCart
func GetOrCreateCart(userId int64) (*Cart, error) {
	cart := &Cart{}

	// Try to get existing active cart
	query := `SELECT id, user_id, status, created_at, updated_at 
	          FROM carts 
	          WHERE user_id=$1 AND status='active'`
	err := db.DB.QueryRow(db.Ctx, query, userId).Scan(&cart.ID, &cart.UserID, &cart.Status, &cart.CreatedAt, &cart.UpdatedAt)

	if err != nil {
		// No active cart found, create new one
		insertQuery := `INSERT INTO carts (user_id, status, created_at) 
		                VALUES ($1, 'active', now()) 
		                RETURNING id, user_id, status, created_at, updated_at`
		err = db.DB.QueryRow(db.Ctx, insertQuery, userId).Scan(&cart.ID, &cart.UserID, &cart.Status, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			return nil, err
		}
	}

	// Load cart items and calculate total
	items, total, err := GetCartItems(cart.ID)
	if err != nil {
		return nil, err
	}

	cart.Items = items
	cart.Total = total

	return cart, nil
}

// Clear removes all items from the cart
// used in: handlers.ClearCart
func (c *Cart) Clear() error {
	query := `DELETE FROM cart_items WHERE cart_id=$1`
	_, err := db.DB.Exec(db.Ctx, query, c.ID)
	if err != nil {
		return err
	}

	// Update cart timestamp
	updateQuery := `UPDATE carts SET updated_at=now() WHERE id=$1`
	_, err = db.DB.Exec(db.Ctx, updateQuery, c.ID)
	return err
}

// Reload refreshes the cart data from database including items and total
// used in: handlers.AddItem, handlers.UpdateItem, handlers.RemoveItem
func (c *Cart) Reload() error {
	items, total, err := GetCartItems(c.ID)
	if err != nil {
		return err
	}
	c.Items = items
	c.Total = total
	return nil
}
