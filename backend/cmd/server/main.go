package main

import (
	"log"

	"auth-backend/internal/config"
	"auth-backend/internal/db"
	"auth-backend/internal/http"
	"auth-backend/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.Load()

	pool, err := db.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	userRepo := user.NewRepository(pool)
	userService := user.NewService(userRepo)

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	http.RegisterRoutes(app, userService)

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}
