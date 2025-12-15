package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/testpilot-ai/llm/adapters"
	"github.com/testpilot-ai/llm/config"
	"github.com/testpilot-ai/llm/handlers"
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

	// Initialize providers
	providerFactory := adapters.NewProviderFactory(
		cfg.OpenAIAPIKey,
		cfg.AnthropicAPIKey,
		cfg.GeminiAPIKey,
		cfg.DefaultProvider,
	)

	log.Printf("Available LLM providers: %v", providerFactory.ListAvailableProviders())
	log.Printf("Default provider: %s", providerFactory.GetDefaultProviderName())

	// Initialize OpenAI provider for embeddings
	openaiProvider := adapters.NewOpenAIProvider(cfg.OpenAIAPIKey)

	// Initialize adapters
	qdrantSearch := adapters.NewQdrantSearchAdapter(cfg.QdrantURL(), "api-knowledge")
	faker := adapters.NewFakerAdapter()
	postgresRepo := adapters.NewPostgresRepository(pool)

	// Initialize handlers
	llmHandler := handlers.NewLLMHandler(
		providerFactory,
		openaiProvider,
		qdrantSearch,
		faker,
		postgresRepo,
	)

	// Setup router
	router := gin.Default()

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

	log.Printf("LLM service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		os.Exit(1)
	}
}

