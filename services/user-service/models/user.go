package models

import (
	"errors"
	"rearatrox/go-ecommerce-backend/pkg/db"
	"rearatrox/go-ecommerce-backend/services/user-service/utils"
)

// User Struct represents a user in the system
type User struct {
	ID           int64   `db:"id" json:"id" swaggerignore:"true"`
	Email        string  `db:"email" json:"email" binding:"required" example:"user@example.com"`
	Password     string  `db:"password" json:"password,omitempty" binding:"required" example:"SecurePass123!" swaggerignore:"true"`
	Role         string  `db:"role" json:"role" example:"user"`
	TokenVersion int     `db:"token_version" json:"-" swaggerignore:"true"`
	FirstName    *string `db:"first_name" json:"firstName,omitempty" example:"Max"`
	LastName     *string `db:"last_name" json:"lastName,omitempty" example:"Mustermann"`
	Phone        *string `db:"phone" json:"phone,omitempty" example:"+49 123 456789"`
}

// GetUsers retrieves all users from the database (admin only)
// used in: handlers.GetUsers
func GetUsers() ([]User, error) {
	query := `SELECT id, email, password, role, token_version, first_name, last_name, phone FROM users`
	rows, err := db.DB.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.TokenVersion, &u.FirstName, &u.LastName, &u.Phone); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// GetUserById retrieves a user by their ID
// used in: handlers.GetUser, handlers.GetMyProfile, handlers.Logout
func GetUserById(id int64) (*User, error) {
	var u User
	query := `SELECT id, email, password, role, token_version, first_name, last_name, phone FROM users WHERE id=$1`
	row := db.DB.QueryRow(db.Ctx, query, id)
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.TokenVersion, &u.FirstName, &u.LastName, &u.Phone); err != nil {
		return nil, err
	}
	return &u, nil
}

// ValidateCredentials checks if the provided password matches the stored hash and loads user data
// used in: handlers.Login
func (u *User) ValidateCredentials() error {
	query := `SELECT password, id, role, token_version FROM users WHERE email=$1`
	row := db.DB.QueryRow(db.Ctx, query, u.Email)
	var hash []byte
	if err := row.Scan(&hash, &u.ID, &u.Role, &u.TokenVersion); err != nil {
		return err
	}

	isPasswordValid := utils.CheckPasswordHash(hash, u.Password)
	if !isPasswordValid {
		return errors.New("credentials invalid")
	}
	return nil
}

// SaveUser creates a new user in the database with hashed password
// used in: handlers.Signup
func (u *User) SaveUser() error {
	query := `INSERT INTO users(email, password, first_name, last_name, phone) VALUES ($1, $2, $3, $4, $5) RETURNING id, role`
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	if err := db.DB.QueryRow(db.Ctx, query, u.Email, hashedPassword, u.FirstName, u.LastName, u.Phone).Scan(&u.ID, &u.Role); err != nil {
		return err
	}
	return nil
}

// UpdateProfile updates the user's profile information (firstName, lastName, phone)
// used in: handlers.UpdateMyProfile
func (u *User) UpdateProfile() error {
	query := `UPDATE users 
	          SET first_name=$1, last_name=$2, phone=$3
	          WHERE id=$4`
	_, err := db.DB.Exec(db.Ctx, query, u.FirstName, u.LastName, u.Phone, u.ID)
	return err
}

// IncrementTokenVersion increments the token version to invalidate all existing JWT tokens
// used in: handlers.Logout
func (u *User) IncrementTokenVersion() error {
	query := `UPDATE users 
	          SET token_version = token_version + 1 
	          WHERE id=$1
	          RETURNING token_version`
	err := db.DB.QueryRow(db.Ctx, query, u.ID).Scan(&u.TokenVersion)
	return err
}
