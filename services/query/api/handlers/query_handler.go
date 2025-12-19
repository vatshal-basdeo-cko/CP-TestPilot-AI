package handlers

import (
	"net/http"
	"strconv"
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

	// Parse pagination
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
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

	// Parse date filters (accept both YYYY-MM-DD and RFC3339 formats)
	if fromDate := c.Query("from_date"); fromDate != "" {
		if t, err := time.Parse("2006-01-02", fromDate); err == nil {
			filters.StartDate = &t
		} else if t, err := time.Parse(time.RFC3339, fromDate); err == nil {
			filters.StartDate = &t
		}
	}

	if toDate := c.Query("to_date"); toDate != "" {
		if t, err := time.Parse("2006-01-02", toDate); err == nil {
			// Set to end of day for inclusive filtering
			endOfDay := t.Add(24*time.Hour - time.Second)
			filters.EndDate = &endOfDay
		} else if t, err := time.Parse(time.RFC3339, toDate); err == nil {
			filters.EndDate = &t
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

// DeleteExecution deletes a test execution record
func (h *QueryHandler) DeleteExecution(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	err = h.historyUseCase.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Execution deleted successfully",
		"id":      idStr,
	})
}

// UpdateValidationResult updates the validation result for an execution
func (h *QueryHandler) UpdateValidationResult(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}

	var body struct {
		ValidationResult map[string]interface{} `json:"validation_result"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.historyUseCase.UpdateValidationResult(c.Request.Context(), id, body.ValidationResult)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Validation result updated successfully",
		"id":      idStr,
	})
}




