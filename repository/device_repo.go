package repository

import (
	"GpsTracker2/database"
	"GpsTracker2/models"
)

// Check if a device exists
func CheckDeviceExist(deviceID int) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Device{}).Where("device_id = ?", deviceID).Count(&count).Error
	return count > 0, err
}

// Create a new device
func CreateDevice(device *models.Device) error {
	return database.DB.Create(device).Error
}

// Get all devices
func GetAllDevices(page, limit int) ([]models.Device, int, error) {
	var devices []models.Device
	var count int64
	totalCount := database.DB.Model(&models.Device{}).Count(&count)
	if totalCount.Error != nil {
		return nil, 0, totalCount.Error
	}

	pageNumber := int((count + int64(limit) - 1) / int64(limit))

	offset := (page - 1) * limit

	err := database.DB.Limit(limit).Offset(offset).Find(&devices).Error
	return devices, pageNumber, err
}

// Get a device by ID
func GetDeviceByID(deviceID int) (models.Device, error) {
	var device models.Device
	err := database.DB.Where("device_id = ?", deviceID).First(&device).Error
	return device, err
}

// Update a device
func UpdateDevice(deviceID int, updatedDevice models.Device) error {
	return database.DB.Model(&models.Device{}).Where("device_id = ?", deviceID).Updates(updatedDevice).Error
}

// Delete a device
func DeleteDevice(deviceID int) error {
	return database.DB.Where("device_id = ?", deviceID).Delete(&models.Device{}).Error
}
