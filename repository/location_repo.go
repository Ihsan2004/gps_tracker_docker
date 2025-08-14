package repository

import (
	"GpsTracker2/database"
	"GpsTracker2/models"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"strconv"
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

// Get location history for a specific device
func GetLocation(deviceID int, page, limit int) ([]models.Location, int, error) {
	var locations []models.Location
	var total int64

	// Count total location
	result := database.DB.Model(&models.Location{}).Count(&total)
	if result.Error != nil {
		return locations, 0, result.Error
	}

	pageNumber := int((total + int64(limit) - 1) / int64(limit))

	offset := (page - 1) * limit
	if err := database.DB.Limit(limit).Offset(offset).Where("device_id = ?", deviceID).Find(&locations).Error; err != nil {
		return nil, 0, err
	}
	return locations, pageNumber, nil
}

// Delete location history for a device
func DeleteLocationHistory(deviceID int) error {
	filter := bson.M{"device_id": deviceID}
	_, err := database.LocationCollection.DeleteMany(context.TODO(), filter)
	return err
}
