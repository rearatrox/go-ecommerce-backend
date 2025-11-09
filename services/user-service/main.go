package main

import (
	"io"
	"log"

	"rearatrox/event-booking-api/pkg/db"
	"rearatrox/event-booking-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

// @title Event Booking API - User-Service
// @version 1.0
// @description API für Event-Verwaltung und Buchung.
// @termsOfService http://swagger.io/terms/

// @contact.name Tim Hauschild
// @contact.url https://webdesign-hauschild.de
// @contact.email info@webdesign-hauschild.de

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:USERSERVICE_PORT
// @BasePath API_PREFIX
func main() {

	if err := logger.InitFromEnv(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync()

	db.InitDB()

	gin.DefaultWriter = io.Discard
	router := gin.Default()

	RegisterRoutes(router)

	router.Run(":8080") // localhost:8080 --> USERSERVICE_PORT mappt dann den Container

}
