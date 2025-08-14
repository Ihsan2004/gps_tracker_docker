package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var LocationCollection *mongo.Collection

// getEnv environment variable ile fallback (MongoDB için)
func getEnvMongo(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// ConnectMongo Docker ortamı için güncellendi
func ConnectMongo() *mongo.Client {
	// Environment variable'dan URI ve database al
	mongoURI := getEnvMongo("MONGODB_URI", "mongodb://localhost:27017")
	mongoDatabase := getEnvMongo("MONGODB_DATABASE", "trackingDB")

	// Set connection options with timeout and retry settings
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(10 * time.Second).
		SetMaxPoolSize(10).
		SetMinPoolSize(1)

	// Connect to MongoDB with retry logic
	var mongoClient *mongo.Client
	var err error

	for i := 0; i < 5; i++ {
		mongoClient, err = mongo.Connect(context.TODO(), clientOptions)
		if err == nil {
			// Ping the database to check connection
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			err = mongoClient.Ping(ctx, nil)
			cancel()

			if err == nil {
				break
			}
		}

		log.Printf("MongoDB connection attempt %d failed: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * 2 * time.Second)
	}

	if err != nil {
		log.Fatal("MongoDB connection failed after retries:", err)
	}

	// Set up the collection - GPS tracking için locationHistory
	LocationCollection = mongoClient.Database(mongoDatabase).Collection("locationHistory")

	log.Printf("Connected to MongoDB successfully at: %s, Database: %s", mongoURI, mongoDatabase)
	return mongoClient
}
