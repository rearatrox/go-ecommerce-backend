package main

import (
	"os"
	"rearatrox/event-booking-api/pkg/logger"
	"rearatrox/event-booking-api/services/user-service/handlers"
	"strings"

	docs "rearatrox/event-booking-api/services/user-service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine) {
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

	port := os.Getenv("USERSERVICE_PORT")
	if port == "" {
		port = "8082"
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = apiPrefix

	api := router.Group(apiPrefix)
	{
		// make sure the swagger UI knows where to fetch the generated spec
		api.GET("/users/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		api.GET("/users", handlers.GetUsers)
		api.GET("/users/:id", handlers.GetUser)

		api.POST("/users/signup", handlers.Signup)
		api.POST("/users/login", handlers.Login)
	}

}
