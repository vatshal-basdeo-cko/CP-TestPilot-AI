package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggingMiddleware logs all requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()

		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// Simple logging (would use structured logger in production)
		println(
			time.Now().Format(time.RFC3339),
			c.Request.Method,
			c.Request.URL.Path,
			statusCode,
			duration.String(),
			requestID,
		)
	}
}




