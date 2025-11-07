package middleware

import (
	"net/http"
	"rearatrox/event-booking-api/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	var token = context.Request.Header.Get("Authorization")
	l := logger.FromContext(context.Request.Context())
	l.Debug("Authenticate required for route")

	if len(token) > 0 {
		parts := strings.Split(token, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			token = parts[1]
		} else {
			l.Error("Invalid Authorization header format")
			context.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid Authorization header format"})
			return
		}
	} else {
		l.Error("Not authorized")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	userId, err := ValidateToken(token)
	if err != nil {
		l.Error("Not authorized")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	l.Debug("Authentication successful")
	context.Set("userId", userId)
	context.Next()
}
