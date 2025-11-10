package main

import (
	"io"
	"log"
	"rearatrox/go-ecommerce-backend/pkg/db"
	"rearatrox/go-ecommerce-backend/pkg/logger"

	"github.com/gin-gonic/gin"
)

// @title Event Booking API - User-Service
// @version 1.0
// @description API für ein E-Commerce Backend
// @termsOfService http://swagger.io/terms/

// @contact.name Tim Hauschild
// @contact.url https://webdesign-hauschild.de
// @contact.email info@webdesign-hauschild.de

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:USERSERVICE_PORT
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

	router.Run(":8080") // localhost:8080 --> USERSERVICE_PORT mappt dann den Container

}

//tmp

// import (
// 	"fmt"
// 	"os"

// 	"golang.org/x/crypto/bcrypt"
// 	"golang.org/x/term"
// )

// func main() {
// 	var pw []byte
// 	if len(os.Args) > 1 {
// 		pw = []byte(os.Args[1])
// 	} else {
// 		fmt.Print("Password: ")
// 		bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
// 		fmt.Println()
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "failed to read password: %v\n", err)
// 			os.Exit(2)
// 		}
// 		pw = bytes
// 	}

// 	cost := bcrypt.DefaultCost // normalerweise 10; du kannst auch 12 setzen: bcrypt.GenerateFromPassword(pw, 12)
// 	hash, err := bcrypt.GenerateFromPassword(pw, cost)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "bcrypt error: %v\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println(string(hash))
// }
