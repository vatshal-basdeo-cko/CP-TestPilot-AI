package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check aggregation
type HealthHandler struct {
	services map[string]string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		services: map[string]string{
			"ingestion":  "http://ingestion:8001/health",
			"llm":        "http://llm:8002/health",
			"execution":  "http://execution:8003/health",
			"validation": "http://validation:8004/health",
			"query":      "http://query:8005/health",
		},
	}
}

// GatewayHealth returns gateway health
func (h *HealthHandler) GatewayHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "gateway",
	})
}

// AllServicesHealth checks health of all services
func (h *HealthHandler) AllServicesHealth(c *gin.Context) {
	statuses := make(map[string]string)

	for name, url := range h.services {
		resp, err := http.Get(url)
		if err != nil {
			statuses[name] = "unhealthy"
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			statuses[name] = "healthy"
		} else {
			statuses[name] = fmt.Sprintf("unhealthy (status: %d)", resp.StatusCode)
		}
	}

	// Determine overall status
	overallHealthy := true
	for _, status := range statuses {
		if status != "healthy" {
			overallHealthy = false
			break
		}
	}

	statusCode := http.StatusOK
	if !overallHealthy {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":   map[bool]string{true: "healthy", false: "degraded"}[overallHealthy],
		"services": statuses,
	})
}

// ProxyHealth proxies health check to a specific service
func (h *HealthHandler) ProxyHealth(c *gin.Context, serviceName string) {
	url, ok := h.services[serviceName]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service not found"})
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "unhealthy",
			"service": serviceName,
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "application/json", body)
}

