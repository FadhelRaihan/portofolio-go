// login.go
package http

import (
	"context"
	"auth-backend/internal/user"

	"github.com/gofiber/fiber/v2"
)

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func loginHandler(userService user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body loginRequest
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
		}

		token, err := userService.Login(context.Background(), body.Email, body.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}

		return c.JSON(fiber.Map{"token": token})
	}
}
