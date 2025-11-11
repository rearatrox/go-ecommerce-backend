package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

type Category struct {
	ID          int64      `db:"id" json:"id" example:"1" swaggerignore:"true"`
	Name        string     `db:"name" json:"name" binding:"required" example:"Elektronik"`
	Slug        string     `db:"slug" json:"slug" binding:"required" example:"elektronik"`
	Description string     `db:"description" json:"description,omitempty" example:"Elektronische Geräte, Zubehör und Gadgets"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt" swaggerignore:"true"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updatedAt,omitempty" swaggerignore:"true"`
}

// InsertCategory creates a new category in the database
// used in: handlers.CreateCategory
func (c *Category) InsertCategory() error {
	query := `INSERT INTO categories (name, slug, description, created_at)
          VALUES ($1, $2, $3, now())
          RETURNING id, created_at`
	if err := db.DB.QueryRow(db.Ctx, query, c.Name, c.Slug, c.Description).Scan(&c.ID, &c.CreatedAt); err != nil {
		return err
	}
	return nil
}

// UpdateCategory updates an existing category's information
// used in: handlers.UpdateCategory
func (c *Category) UpdateCategory() error {
	query := `UPDATE categories
          SET name=$1, slug=$2, description=$3, updated_at=now()
          WHERE id=$4`
	_, err := db.DB.Exec(db.Ctx, query, c.Name, c.Slug, c.Description, c.ID)
	return err
}

// DeleteCategory permanently removes a category from the database
// used in: handlers.DeleteCategoryBySlug
func (c *Category) DeleteCategory() error {
	query := `DELETE FROM categories WHERE id=$1`
	_, err := db.DB.Exec(db.Ctx, query, c.ID)
	return err
}

// GetCategories retrieves all categories from the database ordered by name
// used in: handlers.GetCategories
func GetCategories() ([]Category, error) {
	query := `SELECT id, name, slug, description, created_at, updated_at FROM categories ORDER BY name`
	rows, err := db.DB.Query(db.Ctx, query)
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

// GetCategoryByID retrieves a category by its numeric ID
// used in: handlers.GetCategoryByID
func GetCategoryByID(id int64) (*Category, error) {
	var c Category
	query := `SELECT id, name, slug, description, created_at, updated_at FROM categories WHERE id=$1`
	row := db.DB.QueryRow(db.Ctx, query, id)
	if err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

// GetCategoryBySlug retrieves a category by its URL-friendly slug identifier
// used in: handlers.GetCategoryBySlug, handlers.UpdateCategory, handlers.DeleteCategoryBySlug, handlers.GetProductsByCategory
func GetCategoryBySlug(slug string) (*Category, error) {
	var c Category
	query := `SELECT id, name, slug, description, created_at, updated_at FROM categories WHERE slug=$1`
	row := db.DB.QueryRow(db.Ctx, query, slug)
	if err := row.Scan(&c.ID, &c.Name, &c.Slug, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}
