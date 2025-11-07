package main

import (
	"os"
	"rearatrox/event-booking-api/pkg/logger"
	"rearatrox/event-booking-api/services/user-service/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.Use(logger.GinMiddleware())
	router.NoRoute(func(c *gin.Context) {
		// deinen slog-Logger aus dem Kontext holen
		log := logger.FromContext(c.Request.Context())

		log.Warn("unknown route")

		c.JSON(404, gin.H{"error": "route not found"})
	})

	api := router.Group(os.Getenv("API_PREFIX"))
	{
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUser)

		api.POST("/users/signup", handlers.Signup)
		api.POST("/users/login", handlers.Login)
	}

}
