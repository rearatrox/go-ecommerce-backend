package main

import (
	"os"
	"strings"

	"rearatrox/go-ecommerce-backend/pkg/logger"
	middleware "rearatrox/go-ecommerce-backend/pkg/middleware/auth"
	"rearatrox/go-ecommerce-backend/services/product-service/handlers"

	docs "rearatrox/go-ecommerce-backend/services/product-service/docs"

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

	port := os.Getenv("PRODUCTSERVICE_PORT")
	if port == "" {
		port = "8081"
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = apiPrefix

	api := router.Group(apiPrefix)
	{
		// make sure the swagger UI knows where to fetch the generated spec
		api.GET("/products/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)

		authenticated := api.Group("/")
		{
			authenticated.Use(middleware.Authenticate)
			//authenticated.POST("/products/:id/register", handlers.AddRegistrationForEvent)
			//authenticated.DELETE("/products/:id/delete", handlers.DeleteRegistrationForEvent)

			// admin-only
			admin := authenticated.Group("/admin")
			admin.Use(middleware.Authorize("admin"))
			{
				admin.POST("/products", handlers.CreateProduct)
				admin.PUT("/products/:id", handlers.UpdateProduct)
				admin.DELETE("/products/:id", handlers.DeleteProduct)
			}
		}

	}

}
