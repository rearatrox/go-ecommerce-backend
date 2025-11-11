package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

type Product struct {
	ID          int64      `db:"id" json:"id" swaggerignore:"true"`
	SKU         string     `db:"sku" json:"sku" binding:"required" example:"LAPTOP-001"`
	Name        string     `db:"name" json:"name" binding:"required" example:"Gaming Laptop XPS 15"`
	Description string     `db:"description" json:"description,omitempty" example:"High-performance gaming laptop with RTX 4070"`
	PriceCents  int        `db:"price_cents" json:"priceCents" binding:"required" example:"149999"`
	Currency    string     `db:"currency" json:"currency" example:"EUR"`
	StockQty    int        `db:"stock_qty" json:"stockQty" example:"25"`
	Status      string     `db:"status" json:"status" example:"active"`
	ImageURL    string     `db:"image_url" json:"imageUrl" example:"https://example.com/images/laptop.jpg"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt" swaggerignore:"true"`
	CreatorID   int64      `db:"creator_id" json:"creator_id" swaggerignore:"true"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty" swaggerignore:"true"`
	UpdatorID   int64      `db:"updator_id" json:"updator_id" swaggerignore:"true"`
}

// InsertProduct creates a new product in the database
// used in: handlers.CreateProduct
func (p *Product) InsertProduct() error {
	query := `INSERT INTO products (sku,name,description,price_cents,stock_qty,image_url,creator_id, created_at)
          VALUES ($1,$2,$3,$4,$5,$6,$7, now())
          RETURNING id, status, currency, created_at`
	if err := db.DB.QueryRow(db.Ctx, query, p.SKU, p.Name, p.Description, p.PriceCents,
		p.StockQty, p.ImageURL, p.CreatorID).Scan(&p.ID, &p.Status, &p.Currency, &p.CreatedAt); err != nil {
		return err
	}
	return nil
}

// UpdateProduct updates an existing product's information
// used in: handlers.UpdateProduct
func (p *Product) UpdateProduct() error {
	query := `UPDATE products
          SET name=$1, description=$2, price_cents=$3, currency=$4, stock_qty=$5, status=$6, image_url=$7, updator_id=$8, updated_at=now()
          WHERE sku=$9`
	_, err := db.DB.Exec(db.Ctx, query, p.Name, p.Description, p.PriceCents, p.Currency, p.StockQty, p.Status, p.ImageURL, p.UpdatorID, p.SKU)
	return err
}

// DeleteProductBySKU permanently removes a product from the database
// used in: handlers.DeleteProductBySKU
func (p *Product) DeleteProductBySKU() error {
	query := `DELETE FROM products WHERE sku=$1`
	_, err := db.DB.Exec(db.Ctx, query, p.SKU)
	return err
}

// DeactivateProductBySKU marks a product as inactive without deleting it
// used in: handlers.DeactivateProductBySKU
func (p *Product) DeactivateProductBySKU() error {
	query := `UPDATE products SET status='inactive' WHERE sku=$1`
	_, err := db.DB.Exec(db.Ctx, query, p.SKU)
	return err
}

// GetProducts retrieves all products from the database
// used in: handlers.GetProducts
func GetProducts() ([]Product, error) {
	query := `SELECT id, sku, name, description, price_cents, currency, stock_qty, status, image_url, creator_id, created_at, updated_at FROM products`
	rows, err := db.DB.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents, &p.Currency, &p.StockQty, &p.Status, &p.ImageURL, &p.CreatorID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// GetProductByID retrieves a product by its numeric ID
// used in: handlers.GetProductByID
func GetProductByID(id int64) (*Product, error) {
	var p Product
	query := `SELECT id, sku, name, description, price_cents, currency, stock_qty, status, image_url, creator_id, created_at, updated_at FROM products WHERE id=$1`
	row := db.DB.QueryRow(db.Ctx, query, id)
	if err := row.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents, &p.Currency, &p.StockQty, &p.Status, &p.ImageURL, &p.CreatorID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

// GetProductBySKU retrieves a product by its SKU identifier
// used in: handlers.GetProductBySKU, handlers.UpdateProduct, handlers.DeactivateProductBySKU, handlers.DeleteProductBySKU, handlers.AddCategoriesToProduct, handlers.RemoveCategoryFromProduct, handlers.GetProductCategories
func GetProductBySKU(sku string) (*Product, error) {
	var p Product
	query := `SELECT id, sku, name, description, price_cents, currency, stock_qty, status, image_url, creator_id, created_at, updated_at FROM products WHERE sku=$1`
	row := db.DB.QueryRow(db.Ctx, query, sku)
	if err := row.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents, &p.Currency, &p.StockQty, &p.Status, &p.ImageURL, &p.CreatorID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return nil, err
	}
	return &p, nil
}

// AddCategories assigns multiple categories to a product
// used in: handlers.CreateProduct, handlers.AddCategoriesToProduct
func (p *Product) AddCategories(categoryIds []int64) error {
	if len(categoryIds) == 0 {
		return nil
	}

	for _, categoryId := range categoryIds {
		query := `INSERT INTO product_categories (product_id, category_id, created_at)
		          VALUES ($1, $2, now())
		          ON CONFLICT (product_id, category_id) DO NOTHING`
		_, err := db.DB.Exec(db.Ctx, query, p.ID, categoryId)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveCategory removes a category assignment from a product
// used in: handlers.RemoveCategoryFromProduct
func (p *Product) RemoveCategory(categoryId int64) error {
	query := `DELETE FROM product_categories WHERE product_id=$1 AND category_id=$2`
	_, err := db.DB.Exec(db.Ctx, query, p.ID, categoryId)
	return err
}

// GetProductCategories retrieves all categories assigned to a specific product
// used in: handlers.GetProductCategories
func GetProductCategories(productId int64) ([]Category, error) {
	query := `SELECT c.id, c.name, c.slug, c.description, c.created_at, c.updated_at 
	          FROM categories c
	          INNER JOIN product_categories pc ON c.id = pc.category_id
	          WHERE pc.product_id = $1
	          ORDER BY c.name`
	rows, err := db.DB.Query(db.Ctx, query, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var c Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// GetProductsByCategory retrieves all products assigned to a specific category
// used in: handlers.GetProductsByCategory
func GetProductsByCategory(categoryId int64) ([]Product, error) {
	query := `SELECT p.id, p.sku, p.name, p.description, p.price_cents, p.currency, p.stock_qty, p.status, p.image_url, p.creator_id, p.created_at, p.updated_at
	          FROM products p
	          INNER JOIN product_categories pc ON p.id = pc.product_id
	          WHERE pc.category_id = $1
	          ORDER BY p.name`
	rows, err := db.DB.Query(db.Ctx, query, categoryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.SKU, &p.Name, &p.Description, &p.PriceCents, &p.Currency, &p.StockQty, &p.Status, &p.ImageURL, &p.CreatorID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

// CheckStockAvailable verifies if sufficient stock is available for a product and returns availability status
// used in: handlers.CheckStock
func CheckStockAvailable(productID int64, quantity int) (bool, int, error) {
	var stockQty int
	var status string
	query := `SELECT stock_qty, status FROM products WHERE id=$1`
	err := db.DB.QueryRow(db.Ctx, query, productID).Scan(&stockQty, &status)
	if err != nil {
		return false, 0, err
	}

	// Check if product is active and has enough stock
	if status != "active" {
		return false, stockQty, nil
	}

	return stockQty >= quantity, stockQty, nil
}

// ReduceStock decreases the stock quantity for a product when an order is confirmed
// used in: handlers.ReduceStock
func ReduceStock(productID int64, quantity int) error {
	query := `UPDATE products 
	          SET stock_qty = stock_qty - $1, updated_at = now()
	          WHERE id = $2 AND stock_qty >= $1`
	result, err := db.DB.Exec(db.Ctx, query, quantity, productID)
	if err != nil {
		return err
	}

	// Check if any rows were affected (no rows = insufficient stock or product not found)
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return &StockError{ProductID: productID, Requested: quantity}
	}

	return nil
}

// StockError represents an insufficient stock error
type StockError struct {
	ProductID int64
	Requested int
}

func (e *StockError) Error() string {
	return "insufficient stock or product not found"
}
