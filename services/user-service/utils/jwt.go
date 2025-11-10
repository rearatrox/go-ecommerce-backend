package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email string, userId int64, role string, tokenVersion int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":        email,
		"userId":       userId,
		"role":         role,
		"tokenVersion": tokenVersion,
		"exp":          time.Now().Add(time.Hour * 2).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
