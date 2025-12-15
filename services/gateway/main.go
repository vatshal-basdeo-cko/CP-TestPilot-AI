package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/testpilot-ai/gateway/handlers"
	"github.com/testpilot-ai/gateway/middleware"
	"github.com/testpilot-ai/gateway/proxy"
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

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(pool)
	healthHandler := handlers.NewHealthHandler()
	serviceProxy := proxy.NewServiceProxy()

	// Setup router
	router := gin.Default()

	// Apply global middleware
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())

	// Health routes (no auth required)
	router.GET("/health", healthHandler.GatewayHealth)
	router.GET("/health/all", healthHandler.AllServicesHealth)

	// Auth routes (no auth required)
	auth := router.Group("/api/v1/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// Protected auth routes
	authProtected := router.Group("/api/v1/auth")
	authProtected.Use(middleware.AuthMiddleware())
	{
		authProtected.GET("/me", authHandler.Me)
	}

	// Protected service proxy routes
	router.Any("/api/v1/ingest/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/llm/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/execute/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/validate/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/history/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/apis/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Starting API Gateway on :%s", port)
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




