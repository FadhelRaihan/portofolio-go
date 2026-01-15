package main

import (
	"log"

	"backend/pkg/startup"
)

// @title           Auth & User Management API
// @version         1.0
// @description     API untuk autentikasi dan manajemen user (CRUD).
// @host            localhost:3000
// @BasePath        /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	app, env, err := startup.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Listen(":" + env.AppPort); err != nil {
		log.Fatal(err)
	}
}
