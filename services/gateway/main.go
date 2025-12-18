package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/testpilot-ai/gateway/handlers"
	"github.com/testpilot-ai/gateway/middleware"
	"github.com/testpilot-ai/gateway/proxy"
	"github.com/testpilot-ai/shared/logger"
)

func main() {
	_ = godotenv.Load()

	// Initialize logger
	logger.Init("gateway")

	// Initialize database
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://testpilot:testpilot@postgres:5432/testpilot?sslmode=disable"
	}

	pool, err := initDatabase(databaseURL)
	if err != nil {
		logger.Err(err).Msg("Failed to connect to database")
		os.Exit(1)
	}
	defer pool.Close()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(pool)
	healthHandler := handlers.NewHealthHandler()
	serviceProxy := proxy.NewServiceProxy()

	// Setup router (use gin.New() to avoid default logger noise)
	router := gin.New()
	router.Use(gin.Recovery())

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

	// User management routes (admin only)
	users := router.Group("/api/v1/users")
	users.Use(middleware.AuthMiddleware())
	{
		users.GET("", authHandler.ListUsers)
		users.POST("", authHandler.CreateUser)
		users.DELETE("/:id", authHandler.DeleteUser)
	}

	// Protected service proxy routes
	// Ingestion service
	router.Any("/api/v1/ingest/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/apis", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/apis/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})

	// LLM service
	router.Any("/api/v1/llm/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/parse", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/construct", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})

	// Execution service
	router.Any("/api/v1/execute", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/execute/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/environments", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/environments/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})

	// Validation service
	router.Any("/api/v1/validate", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/validate/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/rules", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/rules/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})

	// Query service
	router.Any("/api/v1/history", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/history/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/analytics", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})
	router.Any("/api/v1/analytics/*path", middleware.AuthMiddleware(), func(c *gin.Context) {
		serviceProxy.RouteToService(c)
	})

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	logger.Infof("Starting API Gateway on :%s", port)
	if err := router.Run(":" + port); err != nil {
		logger.Err(err).Msg("Failed to start server")
		os.Exit(1)
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

	logger.Info("Successfully connected to database")
	return pool, nil
}




