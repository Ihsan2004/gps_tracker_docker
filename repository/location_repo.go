package repository

import (
	"GpsTracker2/database"
	"GpsTracker2/models"
	"context"
	"encoding/json"
	"log"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Create a location history entry
func CreateLocationHistory(location models.Location) error {
	// write to Redis
	marshalJSONlocation, err := json.Marshal(location)
	if err != nil {
		log.Printf("Failed to marshal location: %v", err)
		return err
	}

	err = database.GetRedisClient().Set(
		context.Background(),
		strconv.Itoa(location.DeviceID),
		marshalJSONlocation, 0).Err()
	if err != nil {
		log.Printf("Failed to write to Redis: %v", err)
	}
	return nil
}

// Get location history for a specific device - MONGODB'DEN ÇEK
func GetLocation(deviceID int, page, limit int) ([]models.Location, int, error) {
	var locations []models.Location

	// MongoDB collection'ından total count al
	filter := bson.M{"device_id": deviceID}
	total, err := database.LocationCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.Printf("Failed to count documents: %v", err)
		return locations, 0, err
	}

	// Pagination hesapla
	pageNumber := int((total + int64(limit) - 1) / int64(limit))
	offset := int64((page - 1) * limit)

	// MongoDB'den veri çek - pagination ile
	findOptions := options.Find().
		SetLimit(int64(limit)).
		SetSkip(offset).
		SetSort(bson.D{{Key: "timestamp", Value: -1}}) // En yeni önce

	cursor, err := database.LocationCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		log.Printf("Failed to find documents: %v", err)
		return locations, 0, err
	}
	defer cursor.Close(context.TODO())

	// Cursor'dan sonuçları decode et
	for cursor.Next(context.TODO()) {
		var location models.Location
		if err := cursor.Decode(&location); err != nil {
			log.Printf("Failed to decode location: %v", err)
			continue
		}
		locations = append(locations, location)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return locations, 0, err
	}

	return locations, pageNumber, nil
}

// Delete location history for a device
func DeleteLocationHistory(deviceID int) error {
	filter := bson.M{"device_id": deviceID}
	_, err := database.LocationCollection.DeleteMany(context.TODO(), filter)
	return err
}
