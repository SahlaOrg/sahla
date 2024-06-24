package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mohamed2394/sahla/api/server"
)

func main() {
	// PostgreSQL connection string
	errV := godotenv.Load()
	if errV != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DB")

	// Initialize the server
	srv, err := server.NewServer(dsn)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}
	defer srv.Close()

	// Start the server
	srv.Start(":8080")
}
