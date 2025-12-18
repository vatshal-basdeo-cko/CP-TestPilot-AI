package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/testpilot-ai/validation/adapters"
	"github.com/testpilot-ai/validation/config"
	"github.com/testpilot-ai/validation/handlers"
	"github.com/testpilot-ai/shared/logger"
)

func main() {
	// Load .env file if exists
	_ = godotenv.Load()

	// Initialize logger
	logger.Init("validation")

	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.LogLevel != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to PostgreSQL
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL())
	if err != nil {
		logger.Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(context.Background()); err != nil {
		logger.Err(err).Msg("Failed to ping database")
		os.Exit(1)
	}
	logger.Info("Connected to PostgreSQL")

	// Initialize adapters
	schemaValidator := adapters.NewJSONSchemaValidator()
	postgresRepo := adapters.NewPostgresRepository(pool)

	// Initialize handlers
	validationHandler := handlers.NewValidationHandler(schemaValidator, postgresRepo)

	// Setup router
	router := gin.Default()

	// Health check
	router.GET("/health", validationHandler.Health)

	// API routes
	api := router.Group("/api/v1")
	{
		// Validation endpoints
		api.POST("/validate", validationHandler.Validate)
		api.POST("/compare", validationHandler.Compare)

		// Rules CRUD
		rules := api.Group("/rules")
		{
			rules.GET("", validationHandler.ListRules)
			rules.POST("", validationHandler.CreateRule)
			rules.PUT("/:id", validationHandler.UpdateRule)
			rules.DELETE("/:id", validationHandler.DeleteRule)
		}
	}

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8004"
	}

	logger.Infof("Validation service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Err(err).Msg("Failed to start server")
		os.Exit(1)
	}
}

