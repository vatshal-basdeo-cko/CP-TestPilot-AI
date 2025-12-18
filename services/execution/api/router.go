package api

import (
	"github.com/gin-gonic/gin"
	"github.com/testpilot-ai/execution/api/handlers"
	"github.com/testpilot-ai/execution/api/middleware"
)

// SetupRouter configures and returns the Gin router
func SetupRouter(handler *handlers.ExecutionHandler) *gin.Engine {
	// Use gin.New() to avoid default logger noise
	router := gin.New()
	router.Use(gin.Recovery())

	// Middleware
	router.Use(middleware.RequestLogger())

	// Health check
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Execution endpoints
		v1.POST("/execute", handler.ExecuteAPICall)

		// Environment management
		environments := v1.Group("/environments")
		{
			environments.GET("", handler.ListEnvironments)
			environments.GET("/:id", handler.GetEnvironment)
			environments.POST("", handler.CreateEnvironment)
			environments.PUT("/:id", handler.UpdateEnvironment)
			environments.DELETE("/:id", handler.DeleteEnvironment)
		}
	}

	return router
}

