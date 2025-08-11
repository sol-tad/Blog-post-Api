package main

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/routers"

	// "github.com/gin-gonic/gin"
	docs "github.com/sol-tad/Blog-post-Api/docs"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
)

// @title Blog API
// @version 1.0
// @description This is the Blog API documentation.
// @host localhost:8080
// @BasePath /api/v1

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
	    // Swagger docs setup
    docs.SwaggerInfo.BasePath = "/api/v1"
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Start the server on the specified port
	if port == "" {
    	port = "10000" // Default for local dev
}
	router.Run(":" + port) 

}