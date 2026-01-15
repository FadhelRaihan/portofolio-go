package config

import "os"

type Env struct {
	DBHost string
	DBPort string
	DBUser string
	DBPass string
	DBName string

	AppPort    string
	JWTSecret  string
	JWTIssuer  string
	SwaggerURL string

	GitHubToken	string
	GitHubUsername	string
}

func LoadEnv() Env {
	return Env{
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5433"),
		DBUser: getEnv("DB_USER", "app-user"),
		DBPass: getEnv("DB_PASS", "securedpass!"),
		DBName: getEnv("DB_NAME", "docker-learning-db"),

		AppPort:    getEnv("APP_PORT", "3000"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key"),
		JWTIssuer:  getEnv("JWT_ISSUER", "auth-app"),
		SwaggerURL: getEnv("SWAGGER_URL", "/swagger/*"),

		GitHubToken: getEnv("GITHUB_TOKEN", "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"),
		GitHubUsername: getEnv("GITHUB_USERNAME", "FadhelRaihan"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
