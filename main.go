package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sol-tad/Blog-post-Api/config"
	"github.com/sol-tad/Blog-post-Api/delivery/routers"
)

func main() {
	fmt.Println("GEMINI_API_KEY =", os.Getenv("GEMINI_API_KEY"))
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config.ConnectDB()
	
	port := os.Getenv("PORT")
	router:=routers.SetupRouter()
	router.Run(port)

}