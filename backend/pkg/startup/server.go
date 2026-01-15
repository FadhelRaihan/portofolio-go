package startup

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/middlewares"

	authModule "backend/internal/modules/auth"
	githubmodule "backend/internal/modules/github"
	profileModule "backend/internal/modules/profile"
	userModule "backend/internal/modules/user"

	"backend/pkg/network"
	"backend/pkg/postgres"

	_ "backend/docs"

	fiberSwagger "github.com/swaggo/fiber-swagger"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewServer() (*fiber.App, config.Env, error) {
	env := config.LoadEnv()

	pool, err := postgres.NewPool(env)
	if err != nil {
		return nil, env, err
	}

	jwtProvider := auth.NewJWTProvider(env.JWTSecret, env.JWTIssuer)

	userRepo := userModule.NewRepository(pool)
	userService := userModule.NewService(userRepo, jwtProvider)
	userController := userModule.NewController(userService)

	authController := authModule.NewController(userService)
	authMw := middlewares.Authentication(jwtProvider)

	githubService := githubmodule.NewService(env.GitHubToken)
	githubController := githubmodule.NewController(githubService, env.GitHubUsername)

	profileRepo := profileModule.NewRepository(pool)
	profileService := profileModule.NewService(profileRepo)
	profileController := profileModule.NewController(profileService)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,http://localhost:5174",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	network.RegisterRoutes(app, network.RouterDeps{
		AuthController:   authController,
		UserController:   userController,
		AuthMiddleware:   authMw,
		GitHubController: githubController,
		ProfileController: profileController,
	})

	return app, env, nil
}
