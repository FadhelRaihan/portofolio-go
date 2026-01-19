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

	GitHubToken    string
	GitHubUsername string

	SupabaseURL string
	SupabaseKey string
}

func LoadEnv() Env {
	return Env{
		DBHost: getEnv("DB_HOST", "aws-1-ap-south-1.pooler.supabase.com"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres.eeqdozgfhvhtfxdlqdyb"),
		DBPass: getEnv("DB_PASS", "Padhel29!#."),
		DBName: getEnv("DB_NAME", "postgres"),

		AppPort:    getEnv("APP_PORT", "3000"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key"),
		JWTIssuer:  getEnv("JWT_ISSUER", "auth-app"),
		SwaggerURL: getEnv("SWAGGER_URL", "/swagger/*"),

		GitHubToken:    getEnv("GITHUB_TOKEN", "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"),
		GitHubUsername: getEnv("GITHUB_USERNAME", "FadhelRaihan"),

		SupabaseURL:    getEnv("SUPABASE_URL", "https://eeqdozgfhvhtfxdlqdyb.supabase.co"),
		SupabaseKey:    getEnv("SUPABASE_KEY", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImVlcWRvemdmaHZodGZ4ZGxxZHliIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc2ODQ4MjM5NywiZXhwIjoyMDg0MDU4Mzk3fQ.TMiXws0c1KgmGuy6IdbLM7odzrSyHMoizkL7lMnxWZQ"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
