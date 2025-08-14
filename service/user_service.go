package service

import (
	"GpsTracker2/models"
	"GpsTracker2/repository"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func CreateUser(user *models.User) error {
	existingUser, _ := repository.GetUserByEmail(user.Email)
	if existingUser.ID != 0 {
		return errors.New("Email already exist")
	}
	return repository.CreateUser(user)
}

func GetAllUser(c *fiber.Ctx, page int, limit int) ([]models.APIUser, int, int64, error) {
	users, page_count, total_count, err := repository.GetAllUser(c, page, limit)
	if err != nil {
		return users, 0, 0, c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get all users",
		})
	}
	return users, page_count, total_count, nil
}

func GetUser(userId int) (models.User, error) {
	user, err := repository.GetUser(userId)
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserDevices(userID, page, limit int) ([]models.Device, int, int64, error) {
	userExist, err := repository.CheckUserExist(userID)
	if err != nil {
		return nil, 0, 0, err
	}
	if userExist == 0 {
		return nil, 0, 0, errors.New("User not found")
	}

	result, page_count, total_count, err := repository.GetUserDevices(userID, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}
	return result, page_count, total_count, nil
}

func UpdateUser(userID int, updatedData models.User) error {
	existingUser, _ := repository.CheckUserExist(userID)
	if existingUser == 0 {
		return errors.New("User does not exist")
	}

	// ✅ Pass updatedData to the repository
	return repository.UpdateUser(userID, updatedData)
}

func DeleteUser(userID int) error {
	// ✅ Check if user exists before deleting
	exists, err := repository.CheckUserExist(userID)
	if err != nil {
		return err
	}
	if exists == 0 {
		return errors.New("User not found")
	}

	// ✅ Delete user
	return repository.DeleteUser(userID)
}
