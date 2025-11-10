package models

import (
	"errors"
	"rearatrox/go-ecommerce-backend/pkg/db"
	"rearatrox/go-ecommerce-backend/services/user-service/utils"
)

type User struct {
	ID        int64   `db:"id" json:"id" swaggerignore:"true"`
	Email     string  `db:"email" json:"email" binding:"required" example:"user@example.com"`
	Password  string  `db:"password" json:"password,omitempty" binding:"required" example:"SecurePass123!" swaggerignore:"true"`
	Role      string  `db:"role" json:"role" example:"user"`
	FirstName *string `db:"first_name" json:"firstName,omitempty" example:"Max"`
	LastName  *string `db:"last_name" json:"lastName,omitempty" example:"Mustermann"`
	Phone     *string `db:"phone" json:"phone,omitempty" example:"+49 123 456789"`
}

func GetUsers() ([]User, error) {
	query := `SELECT id, email, password, role, first_name, last_name, phone FROM users`
	rows, err := db.DB.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.FirstName, &u.LastName, &u.Phone); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserById(id int64) (*User, error) {
	var u User
	query := `SELECT id, email, password, role, first_name, last_name, phone FROM users WHERE id=$1`
	row := db.DB.QueryRow(db.Ctx, query, id)
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role, &u.FirstName, &u.LastName, &u.Phone); err != nil {
		return nil, err
	}
	return &u, nil
}

func (u *User) ValidateCredentials() error {
	query := `SELECT password, id, role FROM users WHERE email=$1`
	row := db.DB.QueryRow(db.Ctx, query, u.Email)
	var hash []byte
	if err := row.Scan(&hash, &u.ID, &u.Role); err != nil {
		return err
	}

	//u.Password kommt aus dem Handler --> Request-Eingabe
	isPasswordValid := utils.CheckPasswordHash(hash, u.Password)
	if !isPasswordValid {
		return errors.New("credentials invalid")
	}
	return nil
}

func (u *User) SaveUser() error {
	query := `INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id, role`
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}

	//u.ID wird im User-Objekt gespeichert --> Das User-Objekt wird im Handler aufgerufen -> Es wird die ID angereichert
	if err := db.DB.QueryRow(db.Ctx, query, u.Email, hashedPassword).Scan(&u.ID, &u.Role); err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateProfile() error {
	query := `UPDATE users 
	          SET first_name=$1, last_name=$2, phone=$3
	          WHERE id=$4`
	_, err := db.DB.Exec(db.Ctx, query, u.FirstName, u.LastName, u.Phone, u.ID)
	return err
}
