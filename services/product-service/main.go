package main

import (
	"io"
	"log"

	"rearatrox/go-ecommerce-backend/pkg/db"
	"rearatrox/go-ecommerce-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// @title Event Booking API - Product-Service
// @version 1.0
// @description API für ein E-Commerce Backend
// @termsOfService http://swagger.io/terms/

// @contact.name Tim Hauschild
// @contact.url https://webdesign-hauschild.de
// @contact.email info@webdesign-hauschild.de

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:EVENTSERVICE_PORT
// @BasePath API_PREFIX

// @securityDefinitions.apikey BearerAuth
// @in          header
// @name        Authorization
func main() {

	if err := logger.InitFromEnv(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	db.InitDB()

	gin.DefaultWriter = io.Discard
	router := gin.Default()

	RegisterRoutes(router)

	router.Run(":8080") // localhost:8080 -> Container-Port wird durch EVENTSERVICE_PORT gesetzt

}
