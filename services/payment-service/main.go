package main

import (
	"io"
	"log"
	"os"

	"rearatrox/go-ecommerce-backend/pkg/db"
	"rearatrox/go-ecommerce-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
)

// @title E-Commerce Backend - Payment-Service
// @version 1.0
// @description API für Zahlungsabwicklung im E-Commerce Backend
// @termsOfService http://swagger.io/terms/

// @contact.name Tim Hauschild
// @contact.url https://webdesign-hauschild.de
// @contact.email info@webdesign-hauschild.de

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:PAYMENTSERVICE_PORT
// @BasePath API_PREFIX

// @securityDefinitions.apikey BearerAuth
// @in          header
// @name        Authorization
func main() {

	if err := logger.InitFromEnv(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	// Initialize Stripe
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_SECRET_KEY environment variable is required")
	}
	stripe.Key = stripeKey

	db.InitDB()

	gin.DefaultWriter = io.Discard
	router := gin.Default()

	RegisterRoutes(router)

	router.Run(":8080") // localhost:8080 -> Container-Port wird durch EVENTSERVICE_PORT gesetzt

}
