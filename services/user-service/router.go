package main

import (
	"os"
	"rearatrox/go-ecommerce-backend/pkg/logger"
	middleware "rearatrox/go-ecommerce-backend/pkg/middleware/auth"
	"rearatrox/go-ecommerce-backend/services/user-service/handlers"
	"strings"

	docs "rearatrox/go-ecommerce-backend/services/user-service/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const DEFAULT_PORT = "8081"

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
		port = DEFAULT_PORT
	}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.BasePath = apiPrefix

	api := router.Group(apiPrefix)
	{
		// Public routes
		api.POST("/auth/signup", handlers.Signup)
		api.POST("/auth/login", handlers.Login)

		// make sure the swagger UI knows where to fetch the generated spec
		api.GET("/users/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		// Authenticated routes
		authenticated := api.Group("/")
		authenticated.Use(middleware.Authenticate)
		{
			// Auth endpoints
			authenticated.POST("/auth/logout", handlers.Logout)

			// Profile endpoints
			authenticated.GET("/users/me", handlers.GetMyProfile)
			authenticated.PUT("/users/me", handlers.UpdateMyProfile)

			// Address endpoints
			authenticated.GET("/users/me/addresses", handlers.GetUserAddresses)
			authenticated.GET("/users/me/addresses/:id", handlers.GetAddressByID)
			authenticated.POST("/users/me/addresses", handlers.CreateAddress)
			authenticated.PUT("/users/me/addresses/:id", handlers.UpdateAddress)
			authenticated.DELETE("/users/me/addresses/:id", handlers.DeleteAddress)

			// admin-only
			admin := authenticated.Group("/admin")
			admin.Use(middleware.Authorize("admin"))
			{
				// User list (maybe restrict to admin later)
				admin.GET("/users", handlers.GetUsers)
				admin.GET("/users/:id", handlers.GetUser)
			}
		}
	}

}
