package profile

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewController(s Service) *Controller {
	return &Controller{service: s}
}

// Public: dipakai portfolio
func (c *Controller) GetPublicProfile(ctx *fiber.Ctx) error {
	p, err := c.service.GetCurrentProfile(context.Background())
	if err != nil {
		log.Println("Error:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch profile",
		})
	}
	return ctx.JSON(p)
}

// Admin: list/get by id opsional, di sini contoh CRUD simple.

func (c *Controller) CreateProfile(ctx *fiber.Ctx) error {
	var input CreateProfileInput
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	p, err := c.service.Create(context.Background(), input)
	if err != nil {
		log.Println("Error:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create profile",
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(p)
}

func (c *Controller) UpdateProfile(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	var input UpdateProfileInput
	if err := ctx.BodyParser(&input); err != nil {
		log.Println("Error:", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid body",
		})
	}

	p, err := c.service.Update(context.Background(), id, input)
	if err != nil {
		log.Println("Error:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update profile",
		})
	}
	return ctx.JSON(p)
}

func (c *Controller) DeleteProfile(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id",
		})
	}

	if err := c.service.Delete(context.Background(), id); err != nil {
		log.Println("Error:", err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete profile",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": "Success Delete"})
}
