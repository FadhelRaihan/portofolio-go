package user

import (
	"backend/internal/models"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service Service
}

func NewController(s Service) *Controller {
	return &Controller{service: s}
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	CreatedAt string `json:"created_at"`
}

type userCreateRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type userPatchRequest struct {
	Email    *string `json:"email"`
	FullName *string `json:"full_name"`
	Password *string `json:"password"`
}

func toUserResponse(u *models.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		FullName:  u.FullName,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

func (c *Controller) ListUsers(ctx *fiber.Ctx) error {
	users, err := c.service.ListUsers(context.Background())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list users"})
	}
	res := make([]userResponse, 0, len(users))
	for i := range users {
		res = append(res, toUserResponse(&users[i]))
	}
	return ctx.JSON(res)
}

func (c *Controller) GetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	u, err := c.service.GetUser(context.Background(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	}
	return ctx.JSON(toUserResponse(u))
}

func (c *Controller) CreateUser(ctx *fiber.Ctx) error {
	var body userCreateRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if body.Email == "" || body.FullName == "" || body.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing fields"})
	}
	u, err := c.service.CreateUser(context.Background(), body.Email, body.FullName, body.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create user"})
	}
	return ctx.Status(fiber.StatusCreated).JSON(toUserResponse(u))
}

func (c *Controller) PatchUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var body userPatchRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if body.Email == nil && body.FullName == nil && body.Password == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no fields to update"})
	}
	u, err := c.service.UpdateUser(context.Background(), id, body.Email, body.FullName, body.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update user"})
	}
	return ctx.JSON(toUserResponse(u))
}

func (c *Controller) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := c.service.DeleteUser(context.Background(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete user"})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"success": "Success Delete"})
}
