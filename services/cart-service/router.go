package main

import (
	"os"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	middleware "rearatrox/go-ecommerce-backend/pkg/middleware/auth"
	"rearatrox/go-ecommerce-backend/services/cart-service/handlers"
	"strings"

	docs "rearatrox/go-ecommerce-backend/services/cart-service/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine) {
	router.Use(logger.GinMiddleware())
	router.NoRoute(func(c *gin.Context) {
		log := logger.FromContext(c.Request.Context())
		log.Warn("unknown route")
		c.JSON(404, gin.H{"error": "route not found"})
	})

	// read API prefix, trim spaces and provide a sensible default
	apiPrefix := strings.TrimSpace(os.Getenv("API_PREFIX"))
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	port := os.Getenv("CARTSERVICE_PORT")
	if port == "" {
		port = "8083"
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = apiPrefix

	api := router.Group(apiPrefix)
	{
		// make sure the swagger UI knows where to fetch the generated spec
		api.GET("/cart/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		// All cart routes require authentication
		authenticated := api.Group("/")
		authenticated.Use(middleware.Authenticate)
		{
			// Cart endpoints
			authenticated.GET("/cart", handlers.GetCart)
			authenticated.DELETE("/cart", handlers.ClearCart)

			// Cart items
			authenticated.POST("/cart/items", handlers.AddItem)
			authenticated.PUT("/cart/items/:productId", handlers.UpdateItem)
			authenticated.DELETE("/cart/items/:productId", handlers.RemoveItem)
		}
	}
}
