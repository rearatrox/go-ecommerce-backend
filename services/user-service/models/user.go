package models

import (
	"errors"
	"rearatrox/event-booking-api/pkg/db"
	"rearatrox/event-booking-api/services/user-service/utils"
)

type User struct {
	ID       int64
	Email    string `binding: "required"`
	Password string `binding: "required"`
}

func GetUsers() ([]User, error) {
	query := `SELECT id, email, password FROM users`
	rows, err := db.DB.Query(db.Ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Password); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetUserById(id int64) (*User, error) {
	var u User
	query := `SELECT id, email, password FROM users WHERE id=$1`
	row := db.DB.QueryRow(db.Ctx, query, id)
	if err := row.Scan(&u.ID, &u.Email, &u.Password); err != nil {
		return nil, err
	}
	return &u, nil
}

func (u *User) ValidateCredentials() error {
	query := `SELECT id, password FROM users WHERE email=$1`
	row := db.DB.QueryRow(db.Ctx, query, u.Email)
	var hash []byte
	if err := row.Scan(&u.ID, &hash); err != nil {
		return err
	}
	isPasswordValid := utils.CheckPasswordHash(hash, u.Password)
	if !isPasswordValid {
		return errors.New("credentials invalid")
	}
	return nil
}

func (u *User) SaveUser() error {
	query := `INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id`
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	if err := db.DB.QueryRow(db.Ctx, query, u.Email, hashedPassword).Scan(&u.ID); err != nil {
		return err
	}
	return nil
}
