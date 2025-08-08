package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/routers"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to MongoDB
	config.ConnectDB()

	// Get server port from environment
	port := os.Getenv("PORT")

	// Initialize Gin router with all routes and middleware
	router := routers.SetupRouter()

	// Start the server on the specified port
	router.Run(port)
}