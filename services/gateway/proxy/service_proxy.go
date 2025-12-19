package proxy

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testpilot-ai/shared/logger"
)

// ServiceProxy handles proxying requests to backend services
type ServiceProxy struct {
	services map[string]string
}

// NewServiceProxy creates a new service proxy
func NewServiceProxy() *ServiceProxy {
	return &ServiceProxy{
		services: map[string]string{
			"ingestion": "http://ingestion:8001",
			"llm":       "http://llm:8002",
			"execution": "http://execution:8003",
			"validation": "http://validation:8004",
			"query":     "http://query:8005",
		},
	}
}

// ProxyRequest forwards a request to a backend service
func (sp *ServiceProxy) ProxyRequest(c *gin.Context, serviceName, path string) {
	requestID, _ := c.Get("request_id")
	requestIDStr, _ := requestID.(string)
	
	baseURL, ok := sp.services[serviceName]
	if !ok {
		logger.WithRequestID(requestIDStr).Error().
			Str("service", serviceName).
			Msg("Service not found for proxying")
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service not found"})
		return
	}

	// Build target URL
	targetURL := baseURL + path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Create request
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		logger.WithRequestID(requestIDStr).Err(err).
			Str("service", serviceName).
			Str("url", targetURL).
			Msg("Failed to create proxy request")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add user_id header from context (set by auth middleware)
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			req.Header.Set("X-User-ID", uid.String())
		}
	}

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.WithRequestID(requestIDStr).Err(err).
			Str("service", serviceName).
			Str("url", targetURL).
			Msg("Failed to contact backend service")
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to contact service"})
		return
	}
	defer resp.Body.Close()

	logger.WithRequestID(requestIDStr).Debug().
		Str("service", serviceName).
		Str("path", path).
		Int("status", resp.StatusCode).
		Msg("Proxied request to backend service")

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Copy response body
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// RouteToService determines which service to route to based on path
func (sp *ServiceProxy) RouteToService(c *gin.Context) {
	path := c.Request.URL.Path

	switch {
	// Ingestion service routes
	case strings.HasPrefix(path, "/api/v1/ingest"):
		sp.ProxyRequest(c, "ingestion", path)
	case strings.HasPrefix(path, "/api/v1/apis"):
		sp.ProxyRequest(c, "ingestion", path)

	// LLM service routes - handle both /api/v1/llm/* and direct paths
	case strings.HasPrefix(path, "/api/v1/llm/"):
		// Rewrite /api/v1/llm/X to /api/v1/X for the LLM service
		newPath := strings.Replace(path, "/api/v1/llm/", "/api/v1/", 1)
		sp.ProxyRequest(c, "llm", newPath)
	case strings.HasPrefix(path, "/api/v1/parse"), strings.HasPrefix(path, "/api/v1/construct"):
		sp.ProxyRequest(c, "llm", path)

	// Execution service routes
	case strings.HasPrefix(path, "/api/v1/execute"):
		sp.ProxyRequest(c, "execution", path)
	case strings.HasPrefix(path, "/api/v1/environments"):
		sp.ProxyRequest(c, "execution", path)

	// Validation service routes
	case strings.HasPrefix(path, "/api/v1/validate"):
		sp.ProxyRequest(c, "validation", path)
	case strings.HasPrefix(path, "/api/v1/rules"):
		sp.ProxyRequest(c, "validation", path)

	// Query service routes
	case strings.HasPrefix(path, "/api/v1/history"):
		sp.ProxyRequest(c, "query", path)
	case strings.HasPrefix(path, "/api/v1/analytics"):
		sp.ProxyRequest(c, "query", path)

	default:
		requestID, _ := c.Get("request_id")
		requestIDStr, _ := requestID.(string)
		logger.WithRequestID(requestIDStr).Warn().
			Str("path", path).
			Msg("Route not found")
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
	}
}




