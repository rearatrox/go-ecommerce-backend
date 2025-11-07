package main

import (
	"os"
	"rearatrox/event-booking-api/pkg/logger"
	middleware "rearatrox/event-booking-api/pkg/middleware/auth"
	"rearatrox/event-booking-api/services/event-service/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	// request-logger middleware (adds request-scoped logger into context)
	router.Use(logger.GinMiddleware())
	router.NoRoute(func(c *gin.Context) {
		// deinen slog-Logger aus dem Kontext holen
		log := logger.FromContext(c.Request.Context())

		log.Warn("unknown route")

		c.JSON(404, gin.H{"error": "route not found"})
	})

	api := router.Group(os.Getenv("API_PREFIX"))
	{
		api.GET("/events", handlers.GetEvents)
		api.GET("/events/:id", handlers.GetEvent)

		authenticated := api.Group("/")
		{
			authenticated.Use(middleware.Authenticate)
			authenticated.POST("/events", handlers.CreateEvent)
			authenticated.PUT("/events/:id", handlers.UpdateEvent)
			authenticated.DELETE("/events/:id", handlers.DeleteEvent)

			authenticated.POST("/events/:id/register", handlers.AddRegistrationForEvent)
			authenticated.DELETE("/events/:id/delete", handlers.DeleteRegistrationForEvent)
		}

	}

}
