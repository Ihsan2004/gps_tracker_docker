package routes

import (
	"GpsTracker2/handler"
	"GpsTracker2/middleware"
	"github.com/gofiber/fiber/v2"
)

func ProtectedUserRoutes(app *fiber.App) {
	v1 := app.Group("v1")

	// Public routes
	v1.Post("/login", middleware.LoginHandler)

	// Protected routes with middleware
	v1.Use(middleware.JWTMiddleware)

	// User endpoints
	v1.Post("/users", middleware.RoleMiddleware(middleware.HasAccess), handler.CreateUserHandler)
	v1.Get("/users", middleware.RoleMiddleware(middleware.HasAccess), handler.GetAllUserHandler)
	v1.Get("/users/:userid", middleware.RoleMiddleware(middleware.HasAccess), handler.GetUserHandler)
	v1.Get("/users/devices/:userid", middleware.RoleMiddleware(middleware.HasAccess), handler.GetUserDevicesHandler)
	v1.Put("/users/:userid", middleware.RoleMiddleware(middleware.HasAccess), handler.UpdateUserHandler)
	v1.Delete("/users/:userid", middleware.RoleMiddleware(middleware.HasAccess), handler.DeleteUserHandler)

}
