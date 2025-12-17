package config

import (
	"os"
)

// Config holds all configuration for the ingestion service
type Config struct {
	ServerPort     string
	PostgresHost   string
	PostgresPort   string
	PostgresDB     string
	PostgresUser   string
	PostgresPass   string
	QdrantHost     string
	QdrantPort     string
	GeminiAPIKey   string
	APIConfigsPath string
	LogLevel       string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		ServerPort:     getEnv("SERVER_PORT", "8001"),
		PostgresHost:   getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:   getEnv("POSTGRES_PORT", "5432"),
		PostgresDB:     getEnv("POSTGRES_DB", "testpilot"),
		PostgresUser:   getEnv("POSTGRES_USER", "testpilot"),
		PostgresPass:   getEnv("POSTGRES_PASSWORD", "changeme_in_production"),
		QdrantHost:     getEnv("QDRANT_HOST", "localhost"),
		QdrantPort:     getEnv("QDRANT_PORT", "6333"),
		GeminiAPIKey:   getEnv("GEMINI_API_KEY", ""),
		APIConfigsPath: getEnv("API_CONFIGS_PATH", "./api_configs"),
		LogLevel:       getEnv("LOG_LEVEL", "INFO"),
	}
}

// DatabaseURL returns the PostgreSQL connection URL
func (c *Config) DatabaseURL() string {
	return "postgres://" + c.PostgresUser + ":" + c.PostgresPass + "@" + c.PostgresHost + ":" + c.PostgresPort + "/" + c.PostgresDB
}

// QdrantURL returns the Qdrant connection URL
func (c *Config) QdrantURL() string {
	return "http://" + c.QdrantHost + ":" + c.QdrantPort
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
