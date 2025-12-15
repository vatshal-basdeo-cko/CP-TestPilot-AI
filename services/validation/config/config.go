package config

import (
	"os"
)

// Config holds all configuration for the validation service
type Config struct {
	ServerPort   string
	PostgresHost string
	PostgresPort string
	PostgresDB   string
	PostgresUser string
	PostgresPass string
	LogLevel     string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort:   getEnv("SERVER_PORT", "8004"),
		PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort: getEnv("POSTGRES_PORT", "5432"),
		PostgresDB:   getEnv("POSTGRES_DB", "testpilot"),
		PostgresUser: getEnv("POSTGRES_USER", "testpilot"),
		PostgresPass: getEnv("POSTGRES_PASSWORD", "changeme_in_production"),
		LogLevel:     getEnv("LOG_LEVEL", "INFO"),
	}
}

// DatabaseURL returns the PostgreSQL connection URL
func (c *Config) DatabaseURL() string {
	return "postgres://" + c.PostgresUser + ":" + c.PostgresPass + "@" + c.PostgresHost + ":" + c.PostgresPort + "/" + c.PostgresDB
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

