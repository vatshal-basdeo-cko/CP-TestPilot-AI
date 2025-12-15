package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testpilot-ai/validation/adapters"
	"github.com/testpilot-ai/validation/domain/entities"
)

// ValidationHandler handles validation-related HTTP requests
type ValidationHandler struct {
	schemaValidator *adapters.JSONSchemaValidator
	postgresRepo    *adapters.PostgresRepository
}

// NewValidationHandler creates a new validation handler
func NewValidationHandler(
	schemaValidator *adapters.JSONSchemaValidator,
	postgresRepo *adapters.PostgresRepository,
) *ValidationHandler {
	return &ValidationHandler{
		schemaValidator: schemaValidator,
		postgresRepo:    postgresRepo,
	}
}

// Health returns service health status
func (h *ValidationHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "validation",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// Validate validates an API response
func (h *ValidationHandler) Validate(c *gin.Context) {
	var req entities.ValidationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	result := &entities.ValidationResult{
		IsValid:     true,
		ValidatedAt: time.Now(),
	}

	// Status code validation
	if req.ExpectedStatus > 0 {
		result.StatusCheck = h.schemaValidator.ValidateStatus(req.StatusCode, req.ExpectedStatus)
		if !result.StatusCheck.IsValid {
			result.IsValid = false
			result.Errors = append(result.Errors, "Status code mismatch")
		}
	}

	// Schema validation
	if req.ExpectedSchema != nil {
		result.SchemaCheck = h.schemaValidator.ValidateSchema(req.Response, req.ExpectedSchema)
		if !result.SchemaCheck.IsValid {
			result.IsValid = false
			result.Errors = append(result.Errors, result.SchemaCheck.Errors...)
		}
	}

	// Apply custom rules if API spec ID is provided
	if req.APISpecID != nil {
		rules, err := h.postgresRepo.GetRulesForAPI(c.Request.Context(), *req.APISpecID)
		if err == nil && len(rules) > 0 {
			result.CustomChecks = h.schemaValidator.ApplyCustomRules(req.Response, rules)
			for _, check := range result.CustomChecks {
				if !check.IsValid {
					result.IsValid = false
					result.Errors = append(result.Errors, check.Message)
				}
			}
		}
	}

	// Compare with previous success if provided
	if req.PreviousSuccess != nil {
		diff := h.schemaValidator.CompareResponses(req.Response, req.PreviousSuccess)
		if diff.HasDifferences {
			result.Warnings = append(result.Warnings, "Response differs from previous successful test")
		}
	}

	c.JSON(http.StatusOK, result)
}

// Compare compares two responses
func (h *ValidationHandler) Compare(c *gin.Context) {
	var req struct {
		Current  map[string]interface{} `json:"current" binding:"required"`
		Previous map[string]interface{} `json:"previous" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	diff := h.schemaValidator.CompareResponses(req.Current, req.Previous)
	c.JSON(http.StatusOK, diff)
}

// ListRules lists all validation rules
func (h *ValidationHandler) ListRules(c *gin.Context) {
	rules, err := h.postgresRepo.GetAllRules(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rules": rules,
		"count": len(rules),
	})
}

// CreateRule creates a new validation rule
func (h *ValidationHandler) CreateRule(c *gin.Context) {
	var rule entities.ValidationRule

	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.postgresRepo.CreateRule(c.Request.Context(), &rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rule"})
		return
	}

	c.JSON(http.StatusCreated, rule)
}

// UpdateRule updates an existing validation rule
func (h *ValidationHandler) UpdateRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	var rule entities.ValidationRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	rule.ID = id
	if err := h.postgresRepo.UpdateRule(c.Request.Context(), &rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update rule"})
		return
	}

	c.JSON(http.StatusOK, rule)
}

// DeleteRule deletes a validation rule
func (h *ValidationHandler) DeleteRule(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule ID"})
		return
	}

	if err := h.postgresRepo.DeleteRule(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete rule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Rule deleted"})
}

