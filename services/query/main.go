package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/testpilot-ai/query/api"
	"github.com/testpilot-ai/query/api/handlers"
	"github.com/testpilot-ai/query/application/usecases"
	"github.com/testpilot-ai/query/infrastructure/adapters"
)

func main() {
	_ = godotenv.Load()

	// Initialize database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://testpilot:testpilot@postgres:5432/testpilot?sslmode=disable"
	}

	pool, err := initDatabase(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Initialize repository
	repo := adapters.NewPostgresQueryRepository(pool)

	// Initialize use cases
	historyUseCase := usecases.NewGetTestHistoryUseCase(repo)
	analyticsUseCase := usecases.NewGetAnalyticsUseCase(repo)

	// Initialize handler
	handler := handlers.NewQueryHandler(historyUseCase, analyticsUseCase)

	// Setup router
	router := api.SetupRouter(handler)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8005"
	}

	log.Printf("Starting Query Service on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return pool, nil
}

