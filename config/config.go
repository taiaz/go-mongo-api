package config

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/url"
	"os"
	"time"
)

var MongoDB *mongo.Database

func ConnectDB() {

	mongoUser := os.Getenv("MONGO_USER")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	mongoHostPort := os.Getenv("MONGO_URI")
	mongoDB := os.Getenv("MONGO_DB")

	escapedPassword := url.QueryEscape(mongoPassword)

	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s", mongoUser, escapedPassword, mongoHostPort)
	log.Println("MongoDB URI:", mongoURI)

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Unable to connect to MongoDB:", err)
	}

	log.Println("Connected to MongoDB")

	MongoDB = client.Database(mongoDB)
}
