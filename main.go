package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamed2394/sahla/pkg/db"
)

func main() {
	errV := godotenv.Load()
	if errV != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB")

	// Connect to the PostgreSQL database
	database, err := db.Connect(dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func() {
		dbSQL, err := database.DB()
		if err != nil {
			log.Fatalf("Error getting db from database: %v", err)
		}
		dbSQL.Close()
	}()

	// AutoMigrate models
	if err := db.AutoMigrateModels(); err != nil {
		log.Fatalf("Error migrating models: %v", err)
	}

	log.Println("Database connected and migrated successfully!")
}
