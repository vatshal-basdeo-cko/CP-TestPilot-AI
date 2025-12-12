package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testpilot-ai/query/application/usecases"
	"github.com/testpilot-ai/query/domain/repositories"
)

// QueryHandler handles query-related HTTP requests
type QueryHandler struct {
	historyUseCase   *usecases.GetTestHistoryUseCase
	analyticsUseCase *usecases.GetAnalyticsUseCase
}

// NewQueryHandler creates a new handler
func NewQueryHandler(
	historyUseCase *usecases.GetTestHistoryUseCase,
	analyticsUseCase *usecases.GetAnalyticsUseCase,
) *QueryHandler {
	return &QueryHandler{
		historyUseCase:   historyUseCase,
		analyticsUseCase: analyticsUseCase,
	}
}

// GetHistory retrieves test execution history
func (h *QueryHandler) GetHistory(c *gin.Context) {
	filters := repositories.Filters{
		Status: c.Query("status"),
		Search: c.Query("search"),
		Limit:  20,
		Offset: 0,
	}

	// Parse UUID filters
	if userID := c.Query("user_id"); userID != "" {
		if id, err := uuid.Parse(userID); err == nil {
			filters.UserID = &id
		}
	}

	if apiSpecID := c.Query("api_spec_id"); apiSpecID != "" {
		if id, err := uuid.Parse(apiSpecID); err == nil {
			filters.APISpecID = &id
		}
	}

	executions, total, err := h.historyUseCase.Execute(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"executions": executions,
		"total":      total,
	})
}

// GetExecution retrieves a single execution
func (h *QueryHandler) GetExecution(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	execution, err := h.historyUseCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "execution not found"})
		return
	}

	c.JSON(http.StatusOK, execution)
}

// GetAnalyticsOverview retrieves overall analytics
func (h *QueryHandler) GetAnalyticsOverview(c *gin.Context) {
	var startDate, endDate *time.Time

	if start := c.Query("start_date"); start != "" {
		if t, err := time.Parse(time.RFC3339, start); err == nil {
			startDate = &t
		}
	}

	if end := c.Query("end_date"); end != "" {
		if t, err := time.Parse(time.RFC3339, end); err == nil {
			endDate = &t
		}
	}

	analytics, err := h.analyticsUseCase.GetOverview(c.Request.Context(), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

// GetAPIAnalytics retrieves per-API analytics
func (h *QueryHandler) GetAPIAnalytics(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid API spec ID"})
		return
	}

	stats, err := h.analyticsUseCase.GetByAPI(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// HealthCheck handles health check requests
func (h *QueryHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "query",
	})
}




