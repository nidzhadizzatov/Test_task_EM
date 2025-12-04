package config

import (
	"os"
)

// Config holds the application configuration values
type Config struct {
	DatabaseURL string
	ServerPort  string
	LogLevel    string
}

// LoadConfig reads configuration from environment variables with sensible defaults
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/subscription_service?sslmode=disable"),
		ServerPort:  getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}
	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
