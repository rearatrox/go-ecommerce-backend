package main

import (
	"os"
	"strings"

	"rearatrox/event-booking-api/pkg/logger"
	middleware "rearatrox/event-booking-api/pkg/middleware/auth"
	"rearatrox/event-booking-api/services/event-service/handlers"

	docs "rearatrox/event-booking-api/services/event-service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// read API prefix, trim spaces and provide a sensible default
	apiPrefix := strings.TrimSpace(os.Getenv("API_PREFIX"))
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	port := os.Getenv("EVENTSERVICE_PORT")
	if port == "" {
		port = "8081"
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = apiPrefix

	api := router.Group(apiPrefix)
	{
		// make sure the swagger UI knows where to fetch the generated spec
		api.GET("/events/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		api.GET("/events", handlers.GetEvents)
		api.GET("/events/:id", handlers.GetEvent)

		authenticated := api.Group("/")
		{
			authenticated.Use(middleware.Authenticate)
			authenticated.POST("/events/:id/register", handlers.AddRegistrationForEvent)
			authenticated.DELETE("/events/:id/delete", handlers.DeleteRegistrationForEvent)

			// admin-only
			admin := authenticated.Group("/admin")
			admin.Use(middleware.Authorize("admin"))
			{
				admin.POST("/events", handlers.CreateEvent)
				admin.PUT("/events/:id", handlers.UpdateEvent)
				admin.DELETE("/events/:id", handlers.DeleteEvent)
			}
		}

	}

}
