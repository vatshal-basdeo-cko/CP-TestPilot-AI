package proxy

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
	baseURL, ok := sp.services[serviceName]
	if !ok {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Make request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to contact service"})
		return
	}
	defer resp.Body.Close()

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
	case strings.HasPrefix(path, "/api/v1/ingest"):
		sp.ProxyRequest(c, "ingestion", path)
	case strings.HasPrefix(path, "/api/v1/parse"), strings.HasPrefix(path, "/api/v1/construct"):
		sp.ProxyRequest(c, "llm", path)
	case strings.HasPrefix(path, "/api/v1/execute"), strings.HasPrefix(path, "/api/v1/environments"):
		sp.ProxyRequest(c, "execution", path)
	case strings.HasPrefix(path, "/api/v1/validate"), strings.HasPrefix(path, "/api/v1/rules"):
		sp.ProxyRequest(c, "validation", path)
	case strings.HasPrefix(path, "/api/v1/history"), strings.HasPrefix(path, "/api/v1/analytics"):
		sp.ProxyRequest(c, "query", path)
	default:
		c.JSON(http.StatusNotFound, gin.H{"error": "Route not found"})
	}
}




