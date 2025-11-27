package config

import (
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	ServerPort     string
	DatabaseURL    string
	DefaultTimeout int
	MaxRetries     int
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8003"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://testpilot:testpilot@postgres:5432/testpilot?sslmode=disable"),
		DefaultTimeout: getEnvInt("DEFAULT_TIMEOUT", 30),
		MaxRetries:     getEnvInt("MAX_RETRIES", 3),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

