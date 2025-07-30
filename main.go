package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/routers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDB()
	
	port := os.Getenv("PORT")
	router:=routers.SetupRouter()
	router.Run(port)

}