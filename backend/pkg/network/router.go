package network

import (
	"backend/internal/modules/auth"
	githubmodule "backend/internal/modules/github"
	profileModule "backend/internal/modules/profile"
	skillsModule "backend/internal/modules/skills"
	"backend/internal/modules/user"
	
	"github.com/gofiber/fiber/v2"
)

type RouterDeps struct {
	AuthController    *auth.Controller
	UserController    *user.Controller
	AuthMiddleware    fiber.Handler
	GitHubController  *githubmodule.Controller
	ProfileController *profileModule.Controller
	SkillsController  *skillsModule.Controller
}

func RegisterRoutes(app *fiber.App, deps RouterDeps) {
	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// =========================
	// Public API (no auth required)
	// =========================
	api := app.Group("/api")

	// Login for Dashboard
	api.Post("/login", deps.AuthController.Login)

	// GitHub routes (public)
	github := api.Group("/github")
	github.Get("/me/repos", deps.GitHubController.GetMyReposHandler)
	github.Get("/:username/repos", deps.GitHubController.ListPublicByUserHandler)
	github.Get("/:username/contributions", deps.GitHubController.ContributionsHandler)

	// Profile routes (public read)
	api.Get("/profile", deps.ProfileController.GetPublicProfile)

	// Skills routes (public read)
	api.Get("/skills", deps.SkillsController.List)
	api.Get("/skills/:id", deps.SkillsController.GetByID)

	// =========================
	// Private API (dashboard only - requires JWT auth)
	// =========================
	
	// Protected routes group
	protected := api.Group("", deps.AuthMiddleware)

	// Protected user routes
	users := protected.Group("/users")
	users.Get("/", deps.UserController.ListUsers)
	users.Get("/:id", deps.UserController.GetUser)
	users.Post("/", deps.UserController.CreateUser)
	users.Patch("/:id", deps.UserController.PatchUser)
	users.Delete("/:id", deps.UserController.DeleteUser)

	// Protected profile routes
	profiles := protected.Group("/profiles")
	profiles.Post("/", deps.ProfileController.CreateProfile)
	profiles.Patch("/:id", deps.ProfileController.UpdateProfile)
	profiles.Delete("/:id", deps.ProfileController.DeleteProfile)

	// Protected skills routes (admin only)
	skills := protected.Group("/skills")
	skills.Post("/", deps.SkillsController.Create)
	skills.Put("/:id", deps.SkillsController.Update)
	skills.Delete("/:id", deps.SkillsController.Delete)
}