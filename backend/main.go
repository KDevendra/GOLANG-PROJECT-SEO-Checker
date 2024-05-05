package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() *mongo.Client {
	// Load .env file to retrieve environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Retrieve MongoDB URI from environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalf("MONGO_URI not found in .env file")
	}

	// Create MongoDB client options using the connection URI
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Verify the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB server: %v", err)
	}
	fmt.Println("Connected to MongoDB successfully!")

	return client
}
func RetrieveData(collection *mongo.Collection) ([]interface{}, error) {
	// Define an empty filter to retrieve all documents
	filter := bson.D{}

	// Use Find method to retrieve documents from the collection
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("error finding documents: %v", err)
	}

	// Close the cursor when function completes
	defer cursor.Close(context.Background())

	// Create a slice to hold the retrieved documents
	var results []interface{}

	// Iterate through the cursor to decode each document
	for cursor.Next(context.Background()) {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			return nil, fmt.Errorf("error decoding document: %v", err)
		}
		results = append(results, document)
	}

	// Check for cursor errors
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return results, nil
}

func main() {
	// Connect to MongoDB
	client := Connect()

	// Specify the database and collection
	collection := client.Database("new_database").Collection("posts")

	// Retrieve data from the collection
	data, err := RetrieveData(collection)
	if err != nil {
		log.Fatalf("Error retrieving data: %v", err)
	}

	// Print the retrieved data
	fmt.Println("Retrieved data from MongoDB collection:")
	for _, doc := range data {
		fmt.Println(doc)
	}
}
