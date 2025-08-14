package routes

import (
	"GpsTracker2/handler"
	"GpsTracker2/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProtectedDeviceRoutes(app *fiber.App) {
	v1 := app.Group("v1")

	// Protected routes with middleware
	v1.Use(middleware.JWTMiddleware)

	// Device endpoints
	v1.Post("/devices/:userid", middleware.RoleMiddleware(middleware.HasAccess), handler.CreateDeviceHandler)
	v1.Get("/devices/:device_id", middleware.RoleMiddleware(middleware.HasAccess), handler.GetDeviceHandler)
	v1.Get("/devices", middleware.RoleMiddleware(middleware.HasAccess), handler.GetAllDevicesHandler)
	v1.Put("/devices/:device_id", middleware.RoleMiddleware(middleware.HasAccess), handler.UpdateDeviceHandler)
	v1.Delete("/devices/:device_id", middleware.RoleMiddleware(middleware.HasAccess), handler.DeleteDeviceHandler)

}
