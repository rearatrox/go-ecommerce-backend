package main

import (
	"io"
	"log"

	"rearatrox/event-booking-api/pkg/logger"
	"rearatrox/event-booking-api/services/user-service/db"

	"github.com/gin-gonic/gin"
)

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
