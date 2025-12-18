package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testpilot-ai/shared/logger"
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
		path := c.Request.URL.Path

		// Skip logging successful health checks to reduce noise
		isHealthCheck := strings.HasPrefix(path, "/health")
		isSuccessful := statusCode >= 200 && statusCode < 300

		if isHealthCheck && isSuccessful {
			// Don't log successful health checks
			return
		}

		// Log all other requests, or failed health checks
		log := logger.WithRequestID(requestID)
		logLevel := log.Info()

		// Use error level for failed health checks or server errors
		if isHealthCheck && !isSuccessful {
			logLevel = log.Error()
		} else if statusCode >= 500 {
			logLevel = log.Error()
		}

		logLevel.
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", statusCode).
			Dur("duration", duration).
			Msg("HTTP request")
	}
}
