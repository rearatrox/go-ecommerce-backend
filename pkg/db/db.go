package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool
var Ctx = context.Background()

func InitDB() {

	var err error

	DB, err = pgxpool.New(Ctx, os.Getenv("DATABASE_URL"))

	if err != nil {
		panic("could not connect to the database")
	}

	// Verify the connection
	if err := DB.Ping(Ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}
	log.Println("successfully connected to database")
}
