package network

import (
	"backend/internal/modules/auth"
	githubmodule "backend/internal/modules/github"
	"backend/internal/modules/user"

	"github.com/gofiber/fiber/v2"
)

type RouterDeps struct {
	AuthController   *auth.Controller
	UserController   *user.Controller
	AuthMiddleware   fiber.Handler
	GitHubController *githubmodule.Controller
}

func RegisterRoutes(app *fiber.App, deps RouterDeps) {
	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// =========================
	// Public API (no auth)
	// =========================

	// Login for Dashboard
	app.Post("/api/login", deps.AuthController.Login)

	// GitHub routes for Portofolio
	app.Get("/api/github/me/repos", deps.GitHubController.GetMyReposHandler)
	app.Get("/api/github/:username/repos", deps.GitHubController.ListPublicByUserHandler)
	app.Get("/api/github/:username/contributions", deps.GitHubController.ContributionsHandler)

	// =========================
	// Private API (dashboard only)
	// =========================

	// All routes under /api require AuthMiddleware (JWT)
	api := app.Group("/api", deps.AuthMiddleware)

	api.Get("/users", deps.UserController.ListUsers)
	api.Get("/users/:id", deps.UserController.GetUser)
	api.Post("/users", deps.UserController.CreateUser)
	api.Patch("/users/:id", deps.UserController.PatchUser)
	api.Delete("/users/:id", deps.UserController.DeleteUser)
}
