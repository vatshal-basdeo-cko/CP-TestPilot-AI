package api

import (
	"github.com/gin-gonic/gin"
	"github.com/testpilot-ai/query/api/handlers"
)

// SetupRouter configures the Gin router
func SetupRouter(handler *handlers.QueryHandler) *gin.Engine {
	router := gin.Default()

	// Health check
	router.GET("/health", handler.HealthCheck)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// History endpoints
		v1.GET("/history", handler.GetHistory)
		v1.GET("/history/:id", handler.GetExecution)

		// Analytics endpoints
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/overview", handler.GetAnalyticsOverview)
			analytics.GET("/by-api/:id", handler.GetAPIAnalytics)
		}
	}

	return router
}

