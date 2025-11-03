package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load lädt die .env-Datei aus dem Projektroot.
// Sie kann aus jedem Service aufgerufen werden (einmal pro main.go reicht).
func Load() {
	if err := godotenv.Load("../../.env"); err == nil {
		log.Println("✅ .env-Datei geladen")
	} else {
		if os.Getenv("JWT_SECRET") != "" {
			log.Println("ℹ️  ENV-Variablen aus System übernommen")
		} else {
			log.Println("⚠️  Keine .env-Datei und keine ENV-Variablen gefunden!")
		}
	}
}
