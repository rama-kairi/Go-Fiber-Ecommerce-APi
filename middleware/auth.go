package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rama-kairi/fiber-api/routes/utils"
)

func IsAuthenticated(c *fiber.Ctx) error {
	token_data := c.Get("Authorization")

	if token_data == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized 1",
		})
	}

	token := strings.Split(token_data, " ")[1]

	claims, err := utils.DecodeToken(token, "access")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	c.Locals("user", claims)

	return c.Next()
}

func IsAuthenticatedRefresh(c *fiber.Ctx) error {
	token_data := c.Get("Authorization")

	if token_data == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	token := strings.Split(token_data, " ")[1]

	_, err := utils.DecodeToken(token, "refresh")
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Next()
}
