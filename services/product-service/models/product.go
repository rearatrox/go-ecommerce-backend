package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

type Product struct {
	ID          int64      `db:"id" json:"id"`
	SKU         string     `db:"sku" json:"sku" binding:"required"`
	Name        string     `db:"name" json:"name" binding:"required"`
	Description string     `db:"description" json:"description,omitempty"`
	PriceCents  int        `db:"price_cents" json:"priceCents" binding:"required"`
	Currency    string     `db:"currency" json:"currency"`
	StockQty    int        `db:"stock_qty" json:"stockQty"`
	Status      string     `db:"status" json:"status"`
	ImageURL    string     `db:"image_url" json:"imageUrl"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	CreatorID   int64      `db:"creator_id" json: "creator_id"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

// products slice removed; models now use database storage.

func (p *Product) InsertProduct() error {
	query := `INSERT INTO products (sku,name,description,price_cents,currency,stock_qty,status,image_url,creator_id)
          VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
          RETURNING id, created_at`
	if err := db.DB.QueryRow(db.Ctx, query, p.SKU, p.Name, p.Description, p.PriceCents,
		p.Currency, p.StockQty, p.Status, p.ImageURL, p.CreatorID).Scan(&p.ID, &p.CreatedAt); err != nil {
		return err
	}
	return nil
}

func (p *Product) UpdateProduct() error {
	query := `UPDATE products
          SET name=$1, description=$2, price_cents=$3, currency=$4, stock_qty=$5, status=$6, image_url=$7, updated_at=now()
          WHERE id=$8`
	_, err := db.DB.Exec(db.Ctx, query, p.Name, p.Description, p.PriceCents, p.Currency, p.StockQty, p.Status, p.ImageURL, p.ID)
	return err
}

func (p *Product) DeleteProductByID() error {
	query := `DELETE FROM products WHERE id=$1`
	_, err := db.DB.Exec(db.Ctx, query, p.ID)
	return err
}

func GetProducts() ([]Product, error) {
	query := `SELECT sku, name, description, price_cents, currency, stock_qty, status, image_url FROM products`
	rows, err := db.DB.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.SKU, &p.Name, &p.Description, &p.PriceCents, &p.Currency, &p.StockQty, &p.Status, &p.ImageURL); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func GetProductByID(id int64) (*Product, error) {
	var event Product
	query := `SELECT id, name, description, location, datetime, creator_id FROM products WHERE id=$1`
	row := db.DB.QueryRow(db.Ctx, query, id)
	if err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.CreatorID); err != nil {
		return nil, err
	}
	return &event, nil
}

func (e Product) Register(userId int64) error {
	query := `INSERT INTO event_registrations(event_id, user_id) VALUES ($1, $2)`
	_, err := db.DB.Exec(db.Ctx, query, e.ID, userId)
	return err
}

func (e Product) DeleteRegistration(userId int64) error {
	query := `DELETE FROM event_registrations WHERE event_id=$1 AND user_id=$2`
	_, err := db.DB.Exec(db.Ctx, query, e.ID, userId)
	return err
}
