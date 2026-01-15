package auth

import (
	authdto "backend/internal/modules/auth/dto"
	"backend/internal/modules/user"
	"context"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	userService user.Service
}

func NewController(us user.Service) *Controller {
	return &Controller{userService: us}
}

func (c *Controller) Login(ctx *fiber.Ctx) error {
	var body authdto.LoginRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	token, u, err := c.userService.Login(context.Background(), body.Email, body.Password)
	if err != nil {
		if err == user.ErrInvalidCredentials {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}
	return ctx.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":        u.ID,
			"email":     u.Email,
			"full_name": u.FullName,
		},
	})
}
