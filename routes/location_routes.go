package routes

import (
	"GpsTracker2/handler"
	"GpsTracker2/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProtectedLocationRoutes(app *fiber.App) {
	v1 := app.Group("v1")

	// Protected routes with middleware
	v1.Use(middleware.JWTMiddleware)

	// Websocket endpoint for real-time location updates
	//v1.Get("/locationUpdates/:device_id", middleware.RoleMiddleware(middleware.HasAccess), handler.LocationWebSocketHandler)

	//Location history endpoints
	v1.Post("/locationHistory/:device_id", middleware.RoleMiddleware(middleware.HasAccess), handler.CreateLocationHistoryHandler) //add queue
	v1.Get("/locationHistory/:device_id", middleware.RoleMiddleware(middleware.HasAccess), handler.GetLocationHistoryHandler)
	v1.Delete("/locationHistory/:device_id", middleware.RoleMiddleware(middleware.HasAccess), handler.DeleteLocationHistoryHandler)
}
