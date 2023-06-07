package helper

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB : This is helper function to connect mongoDB
// If you want to export your function. You must to start upper case function name. Otherwise you won't see your function when you import that on other class.
func ConnectDB() *mongo.Database {
	// config := GetConfiguration()
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("library")

	return collection
}

// Configuration model
type Configuration struct {
	Port             string
	ConnectionString string
}

// GetConfiguration method basically populate configuration information from .env and return Configuration model
// func GetConfiguration() Configuration {
// 	err := godotenv.Load("./.env")

// 	if err != nil {
// 		log.Fatalf("Error loading .env file")
// 	}

// 	configuration := Configuration{
// 		os.Getenv("PORT"),
// 		os.Getenv("CONNECTION_STRING"),
// 	}

// 	return configuration
// }
