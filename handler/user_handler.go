package handler

import (
	"GpsTracker2/models"
	"GpsTracker2/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func CreateUserHandler(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	err := service.CreateUser(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	return c.Status(201).JSON(user)
}

func GetAllUserHandler(c *fiber.Ctx) error {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}
	result, pagenumber, total_count, err := service.GetAllUser(c, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get all users",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data": result,
		"paging": fiber.Map{
			"page_count": pagenumber,
			"page_size":  limit,
			"page":       page,
			"total":      total_count,
		},
	})
}

func GetUserHandler(c *fiber.Ctx) error {
	userId := getId(c.Params("userid"))
	result, err := service.GetUser(userId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
	}
	return c.JSON(result)
}

func GetUserDevicesHandler(c *fiber.Ctx) error {
	// Extract user ID from the request parameter
	userID := getId(c.Params("userid"))

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 3
	}

	// Call the service layer
	devices, page_number, total_count, err := service.GetUserDevices(userID, page, limit)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": devices,
		"paging": fiber.Map{
			"page_count": page_number,
			"page_size":  limit,
			"page":       page,
			"total":      total_count,
		},
	})
}

func UpdateUserHandler(c *fiber.Ctx) error {
	userId := getId(c.Params("userid"))

	//  Parse the request body into a user object
	var updatedUser models.User
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	//  Pass `updatedUser` to the service layer
	err := service.UpdateUser(userId, updatedUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    updatedUser,
	})
}

func DeleteUserHandler(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("userid"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	err = service.DeleteUser(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "User deleted successfully"})
}

func getId(s string) int {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return id
}
