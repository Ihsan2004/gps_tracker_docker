package service

import (
	"GpsTracker2/database"
	"GpsTracker2/models"
	"GpsTracker2/repository"
	"encoding/json"
	"errors"
)

// Marshal model location
func MarshalLocation(location models.Location) ([]byte, error) {
	return json.Marshal(&location)
}

// Create Location History
func CreateLocationHistory(deviceID int, location models.Location) error {
	//  Ensure device exists before sending to queue
	exists, err := repository.CheckDeviceExist(deviceID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("Device not found")
	}

	location.DeviceID = deviceID

	// JSON encode location
	body, err := json.Marshal(location)
	if err != nil {
		return err
	}

	// Publish to RabbitMQ
	err = database.PublishMessage(body)
	if err != nil {
		return err
	}

	return repository.CreateLocationHistory(location)
}

// Get Location History
func GetLocationHistory(deviceID, page, pageSize int) ([]models.Location, int, error) {
	return repository.GetLocation(deviceID, page, pageSize)
}

// Delete Location History
func DeleteLocationHistory(deviceID int) error {
	return repository.DeleteLocationHistory(deviceID)
}
