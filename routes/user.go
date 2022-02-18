package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rama-kairi/fiber-api/database"
	"github.com/rama-kairi/fiber-api/models"
	"github.com/rama-kairi/fiber-api/routes/utils"
)

func ResponseUser(user models.User) map[string]interface{} {
	response := map[string]interface{}{
		"id":         user.ID,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
	}
	return response
}

// CreateUser creates a new user
func Signup(c *fiber.Ctx) error {
	type Signup struct {
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	jsonUser := new(Signup)

	if err := c.BodyParser(&jsonUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Check if the password and confirm password match
	if jsonUser.Password != jsonUser.ConfirmPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Passwords do not match",
		})
	}

	// Checking if the Password is Strong
	if validated, msg := utils.PasswordValidator(jsonUser.Password); !validated {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": msg,
		})
	}

	hashed_password, err := utils.HashPassword(jsonUser.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := models.User{
		FirstName: jsonUser.FirstName,
		LastName:  jsonUser.LastName,
		Email:     jsonUser.Email,
		Password:  hashed_password,
	}

	// Handeling the database errors
	if err := database.Database.Db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(ResponseUser(user))
}

// Login - Login a user
func Login(c *fiber.Ctx) error {
	type Login struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	jsonUser := new(Login)

	if err := c.BodyParser(&jsonUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Checking if the User is Authenticated
	isAuthenticated, user := utils.IsAuthenticated(jsonUser.Username, jsonUser.Password)
	if !isAuthenticated {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	// Getting Access Token
	access_token, err := utils.GenerateToken(user.ID, "access")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Getting Refresh Token
	refresh_token, err := utils.GenerateToken(user.ID, "refresh")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":       user.ID,
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}

// UserMe - Returns the current user
func UserMe(c *fiber.Ctx) error {
	tokenString := strings.Split(c.Get("Authorization"), " ")[1]
	claims, err := utils.DecodeToken(tokenString, "access")
	if err != nil {
		return err
	}

	var user models.User
	userId := int(claims["user_id"].(float64))

	database.Database.Db.First(&user, int(userId))

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseUser(user))
}

// Refresh - Refresh the access token
func Refresh(c *fiber.Ctx) error {
	refreshToken := strings.Split(c.Get("Authorization"), " ")[1]
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Refresh token not found",
		})
	}

	claims, err := utils.DecodeToken(refreshToken, "refresh")
	if err != nil {
		return err
	}

	var user models.User
	userId := int(claims["user_id"].(float64))

	database.Database.Db.First(&user, int(userId))

	// Getting Access Token
	access_token, err := utils.GenerateToken(user.ID, "access")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Getting Refresh Token
	refresh_token, err := utils.GenerateToken(user.ID, "refresh")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id":       user.ID,
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}

// GetAllUsers returns all users
func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User

	database.Database.Db.Find(&users)

	responseUsers := make([]map[string]interface{}, len(users))

	for i, user := range users {
		responseUsers[i] = ResponseUser(user)
	}

	return c.Status(fiber.StatusOK).JSON(responseUsers)
}

// GetUser returns a user by id
func GetUser(c *fiber.Ctx) error {
	var user models.User

	id := c.Params("id")

	database.Database.Db.First(&user, id)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found with id " + id,
		})
	}

	return c.Status(fiber.StatusOK).JSON(ResponseUser(user))
}

// UpdateUser updates a user by id
func UpdateUser(c *fiber.Ctx) error {
	var user models.User

	id := c.Params("id")

	database.Database.Db.First(&user, id)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found with id " + id,
		})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	database.Database.Db.Save(&user)

	return c.Status(fiber.StatusOK).JSON(ResponseUser(user))
}

// DeleteUser deletes a user by id
func DeleteUser(c *fiber.Ctx) error {
	var user models.User

	id := c.Params("id")

	database.Database.Db.First(&user, id)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found with id " + id,
		})
	}

	database.Database.Db.Delete(&user)

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}
