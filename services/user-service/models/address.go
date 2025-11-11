package models

import (
	"time"

	"rearatrox/go-ecommerce-backend/pkg/db"
)

// Address Struct represents a user's address
type Address struct {
	ID         int64      `db:"id" json:"id" swaggerignore:"true"`
	UserID     int64      `db:"user_id" json:"userId" swaggerignore:"true"`
	FullName   string     `db:"full_name" json:"fullName" binding:"required" example:"Max Mustermann"`
	Street     string     `db:"street" json:"street" binding:"required" example:"Musterstraße 123"`
	PostalCode string     `db:"postal_code" json:"postalCode" binding:"required" example:"12345"`
	City       string     `db:"city" json:"city" binding:"required" example:"Berlin"`
	Country    string     `db:"country" json:"country" binding:"required" example:"Germany"`
	Type       string     `db:"type" json:"type" binding:"required" example:"shipping"`
	IsDefault  bool       `db:"is_default" json:"isDefault" example:"true"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt" swaggerignore:"true"`
	UpdatedAt  *time.Time `db:"updated_at" json:"updatedAt,omitempty" swaggerignore:"true"`
}

// GetUserAddresses retrieves all addresses for a specific user, ordered by default status and creation date
// used in: handlers.GetUserAddresses
func GetUserAddresses(userId int64) ([]Address, error) {
	query := `SELECT id, user_id, full_name, street, postal_code, city, country, type, is_default, created_at, updated_at 
	          FROM addresses 
	          WHERE user_id=$1 
	          ORDER BY is_default DESC, created_at DESC`
	rows, err := db.DB.Query(db.Ctx, query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var a Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.FullName, &a.Street, &a.PostalCode, &a.City, &a.Country, &a.Type, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

// GetAddressByID retrieves a specific address by ID and user ID to ensure users can only access their own addresses
// used in: handlers.GetAddressByID, handlers.UpdateAddress, handlers.DeleteAddress
func GetAddressByID(addressId int64, userId int64) (*Address, error) {
	var a Address
	query := `SELECT id, user_id, full_name, street, postal_code, city, country, type, is_default, created_at, updated_at 
	          FROM addresses 
	          WHERE id=$1 AND user_id=$2`
	row := db.DB.QueryRow(db.Ctx, query, addressId, userId)
	if err := row.Scan(&a.ID, &a.UserID, &a.FullName, &a.Street, &a.PostalCode, &a.City, &a.Country, &a.Type, &a.IsDefault, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

// InsertAddress creates a new address for a user and unsets other default addresses of the same type if needed
// used in: handlers.CreateAddress
func (a *Address) InsertAddress() error {
	// If this is set as default, unset other defaults of the same type
	if a.IsDefault {
		_, err := db.DB.Exec(db.Ctx,
			`UPDATE addresses SET is_default=false WHERE user_id=$1 AND type=$2`,
			a.UserID, a.Type)
		if err != nil {
			return err
		}
	}

	query := `INSERT INTO addresses (user_id, full_name, street, postal_code, city, country, type, is_default, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
	          RETURNING id, created_at`
	if err := db.DB.QueryRow(db.Ctx, query, a.UserID, a.FullName, a.Street, a.PostalCode, a.City, a.Country, a.Type, a.IsDefault).Scan(&a.ID, &a.CreatedAt); err != nil {
		return err
	}
	return nil
}

// UpdateAddress updates an existing address and unsets other default addresses of the same type if needed
// used in: handlers.UpdateAddress
func (a *Address) UpdateAddress() error {
	// If this is set as default, unset other defaults of the same type
	if a.IsDefault {
		_, err := db.DB.Exec(db.Ctx,
			`UPDATE addresses SET is_default=false WHERE user_id=$1 AND type=$2 AND id!=$3`,
			a.UserID, a.Type, a.ID)
		if err != nil {
			return err
		}
	}

	query := `UPDATE addresses 
	          SET full_name=$1, street=$2, postal_code=$3, city=$4, country=$5, type=$6, is_default=$7, updated_at=now()
	          WHERE id=$8 AND user_id=$9`
	_, err := db.DB.Exec(db.Ctx, query, a.FullName, a.Street, a.PostalCode, a.City, a.Country, a.Type, a.IsDefault, a.ID, a.UserID)
	return err
}

// DeleteAddress deletes an address
// used in: handlers.DeleteAddress
func (a *Address) DeleteAddress() error {
	query := `DELETE FROM addresses WHERE id=$1 AND user_id=$2`
	_, err := db.DB.Exec(db.Ctx, query, a.ID, a.UserID)
	return err
}
