package models

import (
	"errors"
	"rearatrox/go-ecommerce-backend/pkg/db"
	"rearatrox/go-ecommerce-backend/services/user-service/utils"
)

type User struct {
	ID       int64
	Email    string `binding: "required"`
	Password string `binding: "required"`
	Role     string
}

func GetUsers() ([]User, error) {
	query := `SELECT id, email, password, role FROM users`
	rows, err := db.DB.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserById(id int64) (*User, error) {
	var u User
	query := `SELECT id, email, password, role FROM users WHERE id=$1`
	row := db.DB.QueryRow(db.Ctx, query, id)
	if err := row.Scan(&u.ID, &u.Email, &u.Password, &u.Role); err != nil {
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
