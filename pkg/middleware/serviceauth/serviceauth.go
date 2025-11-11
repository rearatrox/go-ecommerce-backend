package serviceauth

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// InternalAuth validates that requests to internal endpoints include the correct secret header
// This prevents external access to service-to-service endpoints
func InternalAuth() gin.HandlerFunc {
	secret := os.Getenv("INTERNAL_API_SECRET")
	if secret == "" {
		panic("INTERNAL_API_SECRET environment variable is required")
	}

	return func(c *gin.Context) {
		headerSecret := c.GetHeader("X-Internal-Secret")

		if headerSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "internal secret required"})
			c.Abort()
			return
		}

		if headerSecret != secret {
			c.JSON(http.StatusForbidden, gin.H{"message": "invalid internal secret"})
			c.Abort()
			return
		}

		c.Next()
	}
}
