package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-mongo-api/config"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"os"
	"time"
)

func GetData(c *gin.Context) {

	collectionName := os.Getenv("MONGO_COLLECTION")
	if collectionName == "" {
		log.Println("Collection name not found in environment variables")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Collection name not found in environment variables"})
		return
	}

	log.Println("Using collection:", collectionName)
	collection := config.MongoDB.Collection(collectionName)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var results []bson.M
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		log.Println("MongoDB Find error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch data from MongoDB", "details": err.Error()})
		return
	}

	if err = cursor.All(ctx, &results); err != nil {
		log.Println("MongoDB Cursor error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to parse data", "details": err.Error()})
		return
	}

	log.Println("Data fetched successfully:", results)
	c.JSON(http.StatusOK, gin.H{"data": results})
}
