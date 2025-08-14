package repository

import (
	"GpsTracker2/database"
	"GpsTracker2/models"
	"github.com/gofiber/fiber/v2"
)

//openapi

func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return user, err
}

// GetAllUser fetches users with pagination
func GetAllUser(c *fiber.Ctx, page int, limit int) ([]models.APIUser, int, int64, error) {
	var users []models.APIUser
	var totalCount int64

	// Count total users
	resultCount := database.DB.Model(&models.User{}).Count(&totalCount)
	if resultCount.Error != nil {
		return nil, 0, 0, resultCount.Error
	}

	// Calculate total pages (round up)
	pageNumber := int((totalCount + int64(limit) - 1) / int64(limit)) // Properly rounds up

	// Ensure page is valid
	if page < 1 {
		page = 1
	}

	// Calculate offset for pagination
	offset := (page - 1) * limit

	// Fetch users with limit, offset, and order
	result := database.DB.Model(&models.User{}).
		Limit(limit).
		Offset(offset).
		Find(&users)

	if result.Error != nil {
		return nil, 0, 0, result.Error
	}

	return users, pageNumber, totalCount, nil
}

func GetUser(userId int) (models.User, error) {
	var user models.User

	result := database.DB.Model(&models.User{}).Where("user_id = ?", userId).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func GetUserDevices(userId int, page int, limit int) ([]models.Device, int, int64, error) {
	var devices []models.Device
	var totalCount int64
	resultCount := database.DB.Model(&models.Device{}).Count(&totalCount)
	if resultCount.Error != nil {
		return nil, 0, 0, resultCount.Error
	}
	pageCount := int((totalCount + int64(limit) - 1) / int64(limit)) // Properly rounds up

	offset := (page - 1) * limit

	result := database.DB.Where("user_id = ?", userId).Limit(limit).Offset(offset).Find(&devices)
	if result.Error != nil {
		return nil, 0, 0, result.Error
	}
	return devices, pageCount, totalCount, nil

}

func UpdateUser(userID int, updatedData models.User) error {
	// ✅ Update only non-empty fields
	return database.DB.Model(&models.User{}).Where("user_id = ?", userID).Updates(updatedData).Error
}

func DeleteUser(userID int) error {
	return database.DB.Where("user_id = ?", userID).Delete(&models.User{}).Error
}

func CheckUserExist(userId int) (int, error) {
	var user models.User
	result := database.DB.Where("user_id = ?", userId).First(&user)
	if result.Error != nil {
		return 0, result.Error
	}
	return userId, nil
}
