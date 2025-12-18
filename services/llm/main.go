package main

import (
	"context"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/testpilot-ai/llm/adapters"
	"github.com/testpilot-ai/llm/config"
	"github.com/testpilot-ai/llm/handlers"
	"github.com/testpilot-ai/shared/logger"
)

func main() {
	// Load .env file if exists
	_ = godotenv.Load()

	// Initialize logger
	logger.Init("llm")

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

	// Initialize providers
	providerFactory := adapters.NewProviderFactory(
		cfg.OpenAIAPIKey,
		cfg.AnthropicAPIKey,
		cfg.GeminiAPIKey,
		cfg.DefaultProvider,
	)

	providers := providerFactory.ListAvailableProviders()
	logger.Infof("LLM providers initialized: [%s], default: %s", 
		strings.Join(providers, ", "), 
		providerFactory.GetDefaultProviderName())

	// Initialize Gemini embedding adapter for RAG
	geminiEmbedding := adapters.NewGeminiEmbeddingAdapter(cfg.GeminiAPIKey)
	if geminiEmbedding.IsAvailable() {
		logger.Info("Gemini embeddings enabled for RAG")
	} else {
		logger.Warn("Gemini embeddings not configured - RAG will not work")
	}

	// Initialize adapters
	qdrantSearch := adapters.NewQdrantSearchAdapter(cfg.QdrantURL(), "api-knowledge")
	faker := adapters.NewFakerAdapter()
	postgresRepo := adapters.NewPostgresRepository(pool)

	// Initialize handlers
	llmHandler := handlers.NewLLMHandler(
		providerFactory,
		geminiEmbedding,
		qdrantSearch,
		faker,
		postgresRepo,
	)

	// Setup router (use gin.New() to avoid default logger noise)
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", llmHandler.Health)

	// API routes
	api := router.Group("/api/v1")
	{
		// LLM endpoints
		api.POST("/parse", llmHandler.ParseRequest)
		api.POST("/construct", llmHandler.ConstructRequest)
		api.POST("/clarify", llmHandler.Clarify)
		api.POST("/generate-data", llmHandler.GenerateData)
		api.POST("/learn", llmHandler.Learn)
		api.GET("/providers", llmHandler.ListProviders)
	}

	// Start server
	port := cfg.ServerPort
	if port == "" {
		port = "8002"
	}

	logger.Infof("LLM service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Err(err).Msg("Failed to start server")
		os.Exit(1)
	}
}
