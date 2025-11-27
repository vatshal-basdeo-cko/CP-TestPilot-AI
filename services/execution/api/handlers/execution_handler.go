package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testpilot-ai/execution/application/usecases"
	"github.com/testpilot-ai/execution/domain/entities"
)

// ExecutionHandler handles execution-related HTTP requests
type ExecutionHandler struct {
	executeUseCase *usecases.ExecuteAPICallUseCase
	envUseCase     *usecases.ManageEnvironmentsUseCase
}

// NewExecutionHandler creates a new execution handler
func NewExecutionHandler(
	executeUseCase *usecases.ExecuteAPICallUseCase,
	envUseCase *usecases.ManageEnvironmentsUseCase,
) *ExecutionHandler {
	return &ExecutionHandler{
		executeUseCase: executeUseCase,
		envUseCase:     envUseCase,
	}
}

// ExecuteAPICall handles API execution requests
func (h *ExecutionHandler) ExecuteAPICall(c *gin.Context) {
	var request entities.APIRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute the API call
	response, err := h.executeUseCase.Execute(c.Request.Context(), &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":    err.Error(),
			"response": response,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListEnvironments handles environment listing
func (h *ExecutionHandler) ListEnvironments(c *gin.Context) {
	environments, err := h.envUseCase.ListEnvironments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"environments": environments,
		"count":        len(environments),
	})
}

// GetEnvironment handles getting a single environment
func (h *ExecutionHandler) GetEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid environment ID"})
		return
	}

	env, err := h.envUseCase.GetEnvironmentByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "environment not found"})
		return
	}

	c.JSON(http.StatusOK, env)
}

// CreateEnvironment handles environment creation
func (h *ExecutionHandler) CreateEnvironment(c *gin.Context) {
	var env entities.Environment
	if err := c.ShouldBindJSON(&env); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.envUseCase.CreateEnvironment(c.Request.Context(), &env); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, env)
}

// UpdateEnvironment handles environment updates
func (h *ExecutionHandler) UpdateEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid environment ID"})
		return
	}

	var env entities.Environment
	if err := c.ShouldBindJSON(&env); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	env.ID = id
	if err := h.envUseCase.UpdateEnvironment(c.Request.Context(), &env); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, env)
}

// DeleteEnvironment handles environment deletion
func (h *ExecutionHandler) DeleteEnvironment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid environment ID"})
		return
	}

	if err := h.envUseCase.DeleteEnvironment(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "environment deleted successfully"})
}

// HealthCheck handles health check requests
func (h *ExecutionHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "execution",
	})
}

