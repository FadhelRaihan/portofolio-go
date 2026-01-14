package config

import (
	"os"
)

type Config struct {
	DBHost	string
	DBPort	string
	DBUser	string
	DBPass	string
	DBName	string
}

func Load() Config {
	return Config{
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5433"),
		DBUser: getEnv("DB_USER", "authuser"),
		DBPass: getEnv("DB_PASS", "auth123"),
		DBName: getEnv("DB_NAME", "docker-learning-db"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}