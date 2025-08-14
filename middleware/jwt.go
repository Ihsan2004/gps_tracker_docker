package middleware

import (
	"GpsTracker2/database"
	"GpsTracker2/models"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var secretKey = []byte("3ccf9f91d8d93a5fb7f4fdd7e1234567890abcdef1234567890abcdef123456")

// Replace with a secure key

type Role struct {
	Role_id   int    `json:"role_id" gorm:"column:role_id"`
	Role_name string `json:"role_name" gorm:"column:role_name"`
}

type Claims struct {
	Username string `json:"Ihsan"`
	Id       int    `json:"id"`
	jwt.RegisteredClaims
}

func createToken(username string, id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Id:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyToken(tokenString string) (*Claims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parsing error: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func LoginHandler(c *fiber.Ctx) error {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user models.User
	if err := database.DB.First(&user, "email = ?", loginData.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	if user.Password != loginData.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	token, err := createToken(user.Name, user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func LoginHandlerHTTP(w http.ResponseWriter, r *http.Request) {
	tokenStr := r.Header.Get("Authorization")
	if tokenStr == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		return
	}

	claims, err := verifyToken(tokenStr)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond with user information if token is valid
	fmt.Fprintf(w, "Hello, %s! Your token is valid.", claims.Username)
}

func JWTMiddleware(c *fiber.Ctx) error {
	if c.Path() == "/v1/login" {
		return c.Next() // Skip JWT validation for /login
	}

	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Authorization header missing",
		})
	}

	// Remove "Bearer " prefix
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	_, err := verifyToken(tokenString)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Check if the user has access to the requested route
	//if ok, message := hasAccess(c, claims.Id); !ok {
	//	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	//		"error": message,
	//	})
	//}

	return c.Next()
}

func RoleMiddleware(hasAccess func(c *fiber.Ctx, userId int) (bool, string)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract Authorization header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header missing",
			})
		}

		// Remove "Bearer " prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Verify the token
		claims, err := verifyToken(tokenString)
		if err != nil || claims == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Extract user ID from claims
		userId := claims.Id

		// Check access using hasAccess function
		if ok, message := hasAccess(c, userId); !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": message,
			})
		}

		// Grant access to the next handler
		return c.Next()
	}
}

func HasAccess(c *fiber.Ctx, userId int) (bool, string) {
	var user models.User

	// Extract and verify the token
	tokenString := c.Get("Authorization")
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims, err := verifyToken(tokenString)
	if err != nil {
		return false, "Invalid or expired token"
	}

	// Fetch the user by ID
	if err := database.DB.First(&user, "user_id = ?", userId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "User not found"
		}
		return false, "Error fetching user"
	}

	// Role-based access logic
	if user.Role_id == 2 { // Regular user
		routeId := getId(c.Params("userid")) // Extract User ID from route params
		if routeId == 0 {                    // If User ID is not in the route, check for Device ID
			deviceId := getId(c.Params("device_id")) // Extract Device ID from route params
			// Fetch the Device record to get the User ID
			var device models.Device
			if err := database.DB.First(&device, "device_id = ?", deviceId).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return false, "Device not found"
				}
				return false, "Error fetching device"
			}

			routeId = device.UserId // Assign User ID from the Device
		}

		fmt.Println("Route ID:", routeId)

		// Check if the token's User ID matches the Route's User ID
		if claims.Id == routeId {
			return true, "User has access to this route"
		}
		return false, "Access denied"
	}

	// Admin role or other roles logic
	if user.Role_id == 1 { // Admin
		return true, "User is an admin"
	}

	return false, "Access denied"
}

func getId(s string) int {
	id, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return id
}
