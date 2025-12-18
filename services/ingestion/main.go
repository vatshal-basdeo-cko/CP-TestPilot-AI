package main

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/testpilot-ai/ingestion/adapters"
	"github.com/testpilot-ai/ingestion/config"
	"github.com/testpilot-ai/ingestion/handlers"
	"github.com/testpilot-ai/shared/logger"
)

func main() {
	// Load .env file if exists
	_ = godotenv.Load()

	// Initialize logger
	logger.Init("ingestion")

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
	fileParser := adapters.NewFileParser()
	postmanParser := adapters.NewPostmanParser()
	embeddingService := adapters.NewEmbeddingService(cfg.GeminiAPIKey)
	qdrantAdapter := adapters.NewQdrantAdapter(cfg.QdrantURL(), "api-knowledge")
	postgresRepo := adapters.NewPostgresRepository(pool)

	// Ensure Qdrant collection exists (768 dimensions for Gemini embeddings)
	if err := qdrantAdapter.EnsureCollection(768); err != nil {
		logger.Err(err).Msg("Failed to ensure Qdrant collection")
	}

	if cfg.GeminiAPIKey != "" {
		logger.Info("Gemini embeddings enabled")
	} else {
		logger.Warn("GEMINI_API_KEY not set - embeddings will be zero vectors")
	}

	// Initialize handlers
	ingestionHandler := handlers.NewIngestionHandler(
		fileParser,
		postmanParser,
		embeddingService,
		qdrantAdapter,
		postgresRepo,
	)

	// Setup router (use gin.New() to avoid default logger noise)
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", ingestionHandler.Health)

	// API routes
	api := router.Group("/api/v1")
	{
		// Ingestion endpoints
		ingest := api.Group("/ingest")
		{
			ingest.POST("/file", ingestionHandler.IngestFile)
			ingest.POST("/folder", ingestionHandler.IngestFolder)
			ingest.POST("/postman", ingestionHandler.IngestPostman)
		}

		// Status and listing
		api.GET("/ingest/status", ingestionHandler.GetStatus)
		api.GET("/apis", ingestionHandler.ListAPIs)
	}

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8001"
	}

	logger.Infof("Ingestion service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Err(err).Msg("Failed to start server")
		os.Exit(1)
	}
}
