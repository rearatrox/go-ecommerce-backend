package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

type CartItem struct {
	ID         int64      `db:"id" json:"id" swaggerignore:"true"`
	CartID     int64      `db:"cart_id" json:"cartId" swaggerignore:"true"`
	ProductID  int64      `db:"product_id" json:"productId" example:"1"`
	Quantity   int        `db:"quantity" json:"quantity" example:"2"`
	PriceCents int        `db:"price_cents" json:"priceCents" example:"2999"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt" swaggerignore:"true"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt,omitempty" swaggerignore:"true"`

	// Product details (joined from products table)
	ProductName        string `json:"productName,omitempty" example:"Gaming Laptop"`
	ProductDescription string `json:"productDescription,omitempty" example:"High-performance laptop"`
	ProductImageURL    string `json:"productImageUrl,omitempty" example:"https://example.com/laptop.jpg"`
}

type AddItemRequest struct {
	ProductID int `json:"productId" example:"1" binding:"required"`
	Quantity  int `json:"quantity" example:"2" binding:"required,min=1"`
}

type UpdateItemRequest struct {
	Quantity int `json:"quantity" example:"3" binding:"required,min=1"`
}

// GetCartItems retrieves all items in a cart with product details and calculates total
func GetCartItems(cartId int64) ([]CartItem, int, error) {
	query := `SELECT 
	            ci.id, ci.cart_id, ci.product_id, ci.quantity, ci.price_cents, 
	            ci.created_at, ci.updated_at,
	            p.name, p.description, p.image_url
	          FROM cart_items ci
	          JOIN products p ON ci.product_id = p.id
	          WHERE ci.cart_id=$1
	          ORDER BY ci.created_at DESC`

	rows, err := db.DB.Query(db.Ctx, query, cartId)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []CartItem
	total := 0

	for rows.Next() {
		var item CartItem
		err := rows.Scan(
			&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.PriceCents,
			&item.CreatedAt, &item.UpdatedAt,
			&item.ProductName, &item.ProductDescription, &item.ProductImageURL,
		)
		if err != nil {
			return nil, 0, err
		}

		items = append(items, item)
		total += item.PriceCents * item.Quantity
	}

	if items == nil {
		items = []CartItem{}
	}

	return items, total, nil
}

// AddOrUpdate adds a product to cart or updates quantity if it already exists
func (ci *CartItem) AddOrUpdate() error {
	// Get current product price
	var priceCents int
	err := db.DB.QueryRow(db.Ctx, `SELECT price_cents FROM products WHERE id=$1`, ci.ProductID).Scan(&priceCents)
	if err != nil {
		return err
	}

	// Check if item already exists in cart
	var existingID int64
	var existingQuantity int
	err = db.DB.QueryRow(db.Ctx, `
		SELECT id, quantity FROM cart_items 
		WHERE cart_id=$1 AND product_id=$2
	`, ci.CartID, ci.ProductID).Scan(&existingID, &existingQuantity)

	if err != nil {
		// Item doesn't exist, insert new
		query := `INSERT INTO cart_items (cart_id, product_id, quantity, price_cents, created_at) 
		          VALUES ($1, $2, $3, $4, now())
		          RETURNING id, created_at`
		err = db.DB.QueryRow(db.Ctx, query, ci.CartID, ci.ProductID, ci.Quantity, priceCents).Scan(&ci.ID, &ci.CreatedAt)
		if err != nil {
			return err
		}
		ci.PriceCents = priceCents
	} else {
		// Item exists, update quantity
		newQuantity := existingQuantity + ci.Quantity
		query := `UPDATE cart_items 
		          SET quantity=$1, updated_at=now() 
		          WHERE id=$2
		          RETURNING updated_at`
		err = db.DB.QueryRow(db.Ctx, query, newQuantity, existingID).Scan(&ci.UpdatedAt)
		if err != nil {
			return err
		}
		ci.ID = existingID
		ci.Quantity = newQuantity
		ci.PriceCents = priceCents
	}

	// Update cart timestamp
	_, err = db.DB.Exec(db.Ctx, `UPDATE carts SET updated_at=now() WHERE id=$1`, ci.CartID)
	return err
}

// UpdateQuantity updates the quantity of a cart item
func (ci *CartItem) UpdateQuantity() error {
	query := `UPDATE cart_items 
	          SET quantity=$1, updated_at=now() 
	          WHERE cart_id=$2 AND product_id=$3
	          RETURNING id, updated_at`
	err := db.DB.QueryRow(db.Ctx, query, ci.Quantity, ci.CartID, ci.ProductID).Scan(&ci.ID, &ci.UpdatedAt)
	if err != nil {
		return err
	}

	// Update cart timestamp
	_, err = db.DB.Exec(db.Ctx, `UPDATE carts SET updated_at=now() WHERE id=$1`, ci.CartID)
	return err
}

// Remove removes a cart item
func (ci *CartItem) Remove() error {
	query := `DELETE FROM cart_items 
	          WHERE cart_id=$1 AND product_id=$2`
	_, err := db.DB.Exec(db.Ctx, query, ci.CartID, ci.ProductID)
	if err != nil {
		return err
	}

	// Update cart timestamp
	_, err = db.DB.Exec(db.Ctx, `UPDATE carts SET updated_at=now() WHERE id=$1`, ci.CartID)
	return err
}
