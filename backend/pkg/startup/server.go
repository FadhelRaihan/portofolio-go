package startup

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/middlewares"

	authModule "backend/internal/modules/auth"
	githubmodule "backend/internal/modules/github"
	profileModule "backend/internal/modules/profile"
	skillsModule "backend/internal/modules/skills"
	userModule "backend/internal/modules/user"

	"backend/pkg/network"
	"backend/pkg/postgres"

	_ "backend/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func NewServer() (*fiber.App, config.Env, error) {
	env := config.LoadEnv()

	// Initialize database connection pool
	pool, err := postgres.NewPool(env)
	if err != nil {
		return nil, env, err
	}

	// Initialize JWT provider
	jwtProvider := auth.NewJWTProvider(env.JWTSecret, env.JWTIssuer)

	// User module
	userRepo := userModule.NewRepository(pool)
	userService := userModule.NewService(userRepo, jwtProvider)
	userController := userModule.NewController(userService)

	// Auth module
	authController := authModule.NewController(userService)
	authMw := middlewares.Authentication(jwtProvider)

	// GitHub module
	githubService := githubmodule.NewService(env.GitHubToken)
	githubController := githubmodule.NewController(githubService, env.GitHubUsername)

	// Profile module
	profileRepo := profileModule.NewRepository(pool)
	profileService := profileModule.NewService(profileRepo)
	profileController := profileModule.NewController(profileService)

	// Skills module - FIXED
	skillsRepo := skillsModule.NewRepository(pool)
	skillsService := skillsModule.NewService(skillsRepo)
	skillsController := skillsModule.NewController(
		skillsService,
		env.SupabaseURL,
		env.SupabaseKey,
	)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB for file uploads
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} (${latency})\n",
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://localhost:5174",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Swagger documentation
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register routes
	network.RegisterRoutes(app, network.RouterDeps{
		AuthController:    authController,
		UserController:    userController,
		AuthMiddleware:    authMw,
		GitHubController:  githubController,
		ProfileController: profileController,
		SkillsController:  skillsController,
	})

	return app, env, nil
}
