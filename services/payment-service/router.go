package main

import (
	"os"
	"strings"

	"rearatrox/go-ecommerce-backend/pkg/logger"
	middleware "rearatrox/go-ecommerce-backend/pkg/middleware/auth"
	"rearatrox/go-ecommerce-backend/services/payment-service/handlers"

	docs "rearatrox/go-ecommerce-backend/services/payment-service/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const DEFAULT_PORT = "8085"

func RegisterRoutes(router *gin.Engine) {
	// CORS Middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// request-logger middleware (adds request-scoped logger into context)
	router.Use(logger.GinMiddleware())
	router.NoRoute(func(c *gin.Context) {
		// logger from context
		log := logger.FromContext(c.Request.Context())

		log.Warn("unknown route")

		c.JSON(404, gin.H{"error": "route not found"})
	})

	// read API prefix, trim spaces and provide a sensible default
	apiPrefix := strings.TrimSpace(os.Getenv("API_PREFIX"))
	if apiPrefix == "" {
		apiPrefix = "/api/v1"
	}

	port := os.Getenv("PAYMENTSERVICE_PORT")
	if port == "" {
		port = DEFAULT_PORT
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = apiPrefix

	api := router.Group(apiPrefix)
	{
		// make sure the swagger UI knows where to fetch the generated spec
		api.GET("/payments/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		// Webhook endpoint (no authentication - verified by Stripe signature)
		api.POST("/webhooks/stripe", handlers.WebhookHandler)

		authenticated := api.Group("/")
		{
			authenticated.Use(middleware.Authenticate)

			// Payment endpoints
			authenticated.POST("/payment-intents", handlers.CreatePaymentIntent)
			authenticated.GET("/payments/:id", handlers.GetPaymentStatus)

			// admin-only
			admin := authenticated.Group("/admin")
			admin.Use(middleware.Authorize("admin"))
			{
				// admin endpoints go here
			}
		}

	}

}
