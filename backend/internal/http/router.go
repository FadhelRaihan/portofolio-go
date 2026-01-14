package http

import (
	"auth-backend/internal/user"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, userService user.Service) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Post("/api/login", loginHandler(userService))
}