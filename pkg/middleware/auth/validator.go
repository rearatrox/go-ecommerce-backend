package middleware

import (
	"context"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ValidateToken(token string, db *pgxpool.Pool, ctx context.Context) (int64, string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return 0, "", errors.New("could not parse token")
	}

	tokenIsValid := parsedToken.Valid

	if !tokenIsValid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid token claims")
	}

	//emailClaim := claims["email"].(string)
	userId := int64(claims["userId"].(float64))
	userRole := claims["role"].(string)
	tokenVersion := int(claims["tokenVersion"].(float64))

	// Validate token version against database
	var dbTokenVersion int
	query := `SELECT token_version FROM users WHERE id=$1`
	err = db.QueryRow(ctx, query, userId).Scan(&dbTokenVersion)
	if err != nil {
		return 0, "", errors.New("user not found")
	}

	if tokenVersion != dbTokenVersion {
		return 0, "", errors.New("token has been revoked")
	}

	return userId, userRole, nil
}
