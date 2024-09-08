package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-mongo-api/config"
	"go-mongo-api/handlers"
	"log"
	"os"
)

func main() {
	envMode := os.Getenv("ENV_MODE")

	if envMode == "local" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		log.Println("Loaded environment variables from .env file")
	} else {
		log.Println("Running in production mode, using environment variables from Kubernetes")
	}

	config.ConnectDB()

	router := gin.Default()

	router.GET("/test-connection", handlers.TestConnection)
	router.GET("/get-data", handlers.GetData)

	router.Run(":8080")
}

//export ENV_MODE=local
//go run main.go
