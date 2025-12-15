package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/testpilot-ai/ingestion/adapters"
	"github.com/testpilot-ai/ingestion/config"
	"github.com/testpilot-ai/ingestion/handlers"
)

func main() {
	// Load .env file if exists
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Set Gin mode
	if cfg.LogLevel != "DEBUG" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to PostgreSQL
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Initialize adapters
	fileParser := adapters.NewFileParser()
	postmanParser := adapters.NewPostmanParser()
	embeddingService := adapters.NewEmbeddingService(cfg.OpenAIAPIKey)
	qdrantAdapter := adapters.NewQdrantAdapter(cfg.QdrantURL(), "api-knowledge")
	postgresRepo := adapters.NewPostgresRepository(pool)

	// Ensure Qdrant collection exists
	if err := qdrantAdapter.EnsureCollection(1536); err != nil {
		log.Printf("Warning: Failed to ensure Qdrant collection: %v", err)
	}

	// Initialize handlers
	ingestionHandler := handlers.NewIngestionHandler(
		fileParser,
		postmanParser,
		embeddingService,
		qdrantAdapter,
		postgresRepo,
	)

	// Setup router
	router := gin.Default()

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

	log.Printf("Ingestion service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}

