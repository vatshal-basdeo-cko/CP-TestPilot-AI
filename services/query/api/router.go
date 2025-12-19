package api

import (
	"github.com/gin-gonic/gin"
	"github.com/testpilot-ai/query/api/handlers"
)

// SetupRouter configures the Gin router
func SetupRouter(handler *handlers.QueryHandler) *gin.Engine {
	// Use gin.New() to avoid default logger noise
	router := gin.New()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// History endpoints
		v1.GET("/history", handler.GetHistory)
		v1.GET("/history/:id", handler.GetExecution)
		v1.DELETE("/history/:id", handler.DeleteExecution)
		v1.PATCH("/history/:id/validation", handler.UpdateValidationResult)

		// Analytics endpoints
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/overview", handler.GetAnalyticsOverview)
			analytics.GET("/by-api/:id", handler.GetAPIAnalytics)
		}
	}

	return router
}




