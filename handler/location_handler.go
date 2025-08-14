package handler

import (
	"GpsTracker2/models"
	"GpsTracker2/service"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"strconv"
	"sync"
)

// new variables for managing websocket connections
var (
	// Bir device ID için birden fazla WebSocket bağlantısı tutmak için
	activeConnections map[string][]*websocket.Conn
	connectionsMutex  sync.RWMutex
)

// init fonksiyonu ile map'i başlat
func init() {
	activeConnections = make(map[string][]*websocket.Conn)
}

// Adding connection function for websocket connections
func addConnection(deviceID string, conn *websocket.Conn) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	// If there is no slice for this device ID, create one
	if activeConnections[deviceID] == nil {
		activeConnections[deviceID] = make([]*websocket.Conn, 0)
	}

	activeConnections[deviceID] = append(activeConnections[deviceID], conn)
	log.Printf("Added connection for device ID %s. Total connections : %d", deviceID, len(activeConnections[deviceID]))
}

// Removing connection function for websocket connections
func removeConnection(deviceID string, conn *websocket.Conn) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	connections, exists := activeConnections[deviceID]
	if !exists {
		return
	}

	// Belirli connection'ı bul ve kaldır
	for i, c := range connections {
		if c == conn {
			activeConnections[deviceID] = append(connections[:i], connections[i+1:]...)
			break
		}
	}

	// Eğer hiç connection kalmadıysa map'ten kaldır
	if len(activeConnections[deviceID]) == 0 {
		delete(activeConnections, deviceID)
	}

	log.Printf("Connection removed for device %s. Remaining connections: %d", deviceID, len(activeConnections[deviceID]))
}

// Broadcast function to send messages to all connections of a specific device
func broadcastToDevice(deviceID string, message interface{}) {
	connectionsMutex.RLock()
	connections, exists := activeConnections[deviceID]
	defer connectionsMutex.RUnlock()

	if !exists || len(connections) == 0 {
		log.Printf("No active connections for device ID %s", deviceID)
		return
	}

}

// Create Location History
func CreateLocationHistoryHandler(c *fiber.Ctx) error {
	deviceID, err := strconv.Atoi(c.Params("device_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid device ID"})
	}

	var location models.Location
	if err := c.BodyParser(&location); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err = service.CreateLocationHistory(deviceID, location)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(location)
}

// Get Location History by Device ID
func GetLocationHistoryHandler(c *fiber.Ctx) error {
	deviceID, err := strconv.Atoi(c.Params("device_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid device ID"})
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 3
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 2
	}
	history, pageNumber, err := service.GetLocationHistory(deviceID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data":       history,
		"total_page": pageNumber,
	})
}

// Delete Location History by Device ID
func DeleteLocationHistoryHandler(c *fiber.Ctx) error {
	deviceID, err := strconv.Atoi(c.Params("device_id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid device ID"})
	}

	err = service.DeleteLocationHistory(deviceID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Location history deleted successfully"})
}

func LocationWebSocketHandler(c *websocket.Conn) {

	// Get device ID from URL parameters
	deviceID := c.Params("device_id")

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}

		// Process the message - parse it as a location update
		var location models.Location
		if err := json.Unmarshal(msg, &location); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Error : Invalid location data"))
			continue
		}

		// Save to database
		deviceIDInt, _ := strconv.Atoi(deviceID)
		if err := service.CreateLocationHistory(deviceIDInt, location); err != nil {
			c.WriteMessage(websocket.TextMessage, []byte("Error : Failed to save location data"))
			continue
		}

		// Send confirmation back to the client
		c.WriteMessage(websocket.TextMessage, []byte("Location data received and saved successfully"))
	}
}
