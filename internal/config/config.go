package config

import "os"

type Config struct {
	Port         string
	DatabaseURL  string
	SessionKey   string
	JWTSecret    string
	Environment  string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/expenses?sslmode=disable"),
		SessionKey:  getEnv("SESSION_KEY", "change-me-in-production"),
		JWTSecret:   getEnv("JWT_SECRET", "change-me-in-production"),
		Environment: getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
