package main

import (
	"os"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	middleware "rearatrox/go-ecommerce-backend/pkg/middleware/auth"
	"rearatrox/go-ecommerce-backend/pkg/middleware/serviceauth"
	"rearatrox/go-ecommerce-backend/services/order-service/handlers"
	"strings"

	docs "rearatrox/go-ecommerce-backend/services/order-service/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const DEFAULT_PORT = "8084"

func RegisterRoutes(router *gin.Engine) {
	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

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

	port := os.Getenv("ORDERSERVICE_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = apiPrefix

	api := router.Group(apiPrefix)
	{
		// make sure the swagger UI knows where to fetch the generated spec
		api.GET("/orders/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		// Internal endpoints (service-to-service communication with secret)
		internal := api.Group("/internal")
		internal.Use(serviceauth.InternalAuth())
		{
			internal.PATCH("/orders/:id/status", handlers.InternalUpdateOrderStatus)
		}

		// All order routes require authentication
		authenticated := api.Group("/")
		authenticated.Use(middleware.Authenticate)
		{
			// Order endpoints
			authenticated.POST("/orders", handlers.CreateOrder)
			authenticated.GET("/orders", handlers.ListOrders)
			authenticated.GET("/orders/:id", handlers.GetOrder)
			authenticated.PATCH("/orders/:id/status", handlers.UpdateOrderStatus)
			authenticated.PATCH("/orders/:id/cancel", handlers.CancelOrder)
		}
	}
}
