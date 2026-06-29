// Package config loads runtime configuration from environment variables.
package config

import "os"

// Config holds the application configuration resolved at startup.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
}

// Load reads configuration from the environment, applying sane defaults for
// local development.
func Load() Config {
	return Config{
		Port:        getenv("PORT", "8080"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/credit_analysis?sslmode=disable"),
		JWTSecret:   getenv("JWT_SECRET", "change-me-in-production"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
