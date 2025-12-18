package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/testpilot-ai/execution/api"
	"github.com/testpilot-ai/execution/api/handlers"
	"github.com/testpilot-ai/execution/application/usecases"
	"github.com/testpilot-ai/execution/infrastructure/adapters"
	"github.com/testpilot-ai/execution/infrastructure/config"
	"github.com/testpilot-ai/shared/logger"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Initialize logger
	logger.Init("execution")

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	pool, err := initDatabase(cfg.DatabaseURL)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize repositories
	executionRepo := adapters.NewPostgresRepository(pool)
	envRepo := adapters.NewEnvironmentRepository(pool)

	// Initialize use cases
	executeUseCase := usecases.NewExecuteAPICallUseCase(executionRepo)
	envUseCase := usecases.NewManageEnvironmentsUseCase(envRepo)

	// Initialize handlers
	handler := handlers.NewExecutionHandler(executeUseCase, envUseCase)

	// Setup router
	router := api.SetupRouter(handler)

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	logger.Infof("Starting Execution Service on %s", serverAddr)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(serverAddr); err != nil {
			logger.Err(err).Msg("Failed to start server")
			os.Exit(1)
		}
	}()

	<-quit
	logger.Info("Shutting down server...")
}

func initDatabase(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	logger.Info("Successfully connected to database")
	return pool, nil
}

