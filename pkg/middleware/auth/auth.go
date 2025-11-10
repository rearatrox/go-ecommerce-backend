package middleware

import (
	"net/http"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	var token = context.Request.Header.Get("Authorization")
	l := logger.FromContext(context.Request.Context())
	l.Debug("Authentication required for route")

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

	userId, userRole, err := ValidateToken(token)
	if err != nil {
		l.Error("Not authorized")
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	l.Debug("Authentication successful")
	context.Set("userId", userId)
	context.Set("userRole", userRole)
	context.Next()
}

const (
	CtxRole = "userRole"
)

func Authorize(allowed ...string) gin.HandlerFunc {
	allow := make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		allow[a] = struct{}{}
	}

	return func(c *gin.Context) {
		v, ok := c.Get(CtxRole)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no role in context"})
			return
		}
		role, ok := v.(string)
		if !ok || role == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid role"})
			return
		}

		// admin passt immer durch
		if role == "admin" {
			c.Next()
			return
		}
		if _, ok := allow[role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient role"})
			return
		}
		c.Next()
	}
}
