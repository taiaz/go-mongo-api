package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-mongo-api/config"
	"go-mongo-api/handlers"
	"log"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize MongoDB connection
	config.ConnectDB()

	// Create a new Gin router
	router := gin.Default()

	// Routes
	router.GET("/test-connection", handlers.TestConnection)
	router.GET("/get-data", handlers.GetData)

	// Start the server
	router.Run(":8080")
}
