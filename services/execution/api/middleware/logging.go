package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestLogger logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		requestID := uuid.New().String()

		// Set request ID
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Log request details
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

