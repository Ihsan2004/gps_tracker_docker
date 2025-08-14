package handler

import (
	"GpsTracker2/database"
	"GpsTracker2/models"
	"GpsTracker2/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

// Create Device
func CreateDeviceHandler(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("userid"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var device models.Device
	if err := c.BodyParser(&device); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err = service.CreateDevice(userID, &device, database.ConnectElastic())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(device)
}

// Get All Devices
func GetAllDevicesHandler(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(c.Query("limit"))
	if err != nil || size < 1 {
		size = 2
	}

	query := c.Query("query")

	// Parse category as comma-separated values, e.g. "1,2"
	categoryStr := c.Query("category")
	var categories []int
	if categoryStr != "" {
		for _, s := range strings.Split(categoryStr, ",") {
			cat, err := strconv.Atoi(strings.TrimSpace(s))
			if err == nil {
				categories = append(categories, cat)
			}
		}
	}

	devices, total_page, err := service.GetAllDevices(page, size, database.ConnectElastic(), query, categories)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       devices,
		"total_page": total_page,
	})
}

// Get Device by ID
func GetDeviceHandler(c *fiber.Ctx) error {
	deviceID, err := strconv.Atoi(c.Params("device_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid device ID"})
	}

	device, err := service.GetDeviceByID(deviceID, database.ConnectElastic())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(device)
}

// Update Device by ID
func UpdateDeviceHandler(c *fiber.Ctx) error {
	deviceID, err := strconv.Atoi(c.Params("device_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid device ID"})
	}

	var updatedDevice models.Device
	if err := c.BodyParser(&updatedDevice); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err = service.UpdateDevice(deviceID, updatedDevice, database.ConnectElastic())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Device updated successfully"})
}

// Delete Device by ID
func DeleteDeviceHandler(c *fiber.Ctx) error {
	deviceID, err := strconv.Atoi(c.Params("device_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid device ID"})
	}

	err = service.DeleteDevice(deviceID, database.ConnectElastic())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Device deleted successfully"})
}
