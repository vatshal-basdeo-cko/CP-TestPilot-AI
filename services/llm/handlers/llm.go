package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testpilot-ai/llm/adapters"
	"github.com/testpilot-ai/llm/domain/entities"
	"github.com/testpilot-ai/llm/prompts"
	"github.com/testpilot-ai/shared/logger"
)

// LLMHandler handles LLM-related HTTP requests
type LLMHandler struct {
	providerFactory *adapters.ProviderFactory
	geminiEmbedding *adapters.GeminiEmbeddingAdapter
	qdrantSearch    *adapters.QdrantSearchAdapter
	faker           *adapters.FakerAdapter
	postgresRepo    *adapters.PostgresRepository
}

// NewLLMHandler creates a new LLM handler
func NewLLMHandler(
	providerFactory *adapters.ProviderFactory,
	geminiEmbedding *adapters.GeminiEmbeddingAdapter,
	qdrantSearch *adapters.QdrantSearchAdapter,
	faker *adapters.FakerAdapter,
	postgresRepo *adapters.PostgresRepository,
) *LLMHandler {
	return &LLMHandler{
		providerFactory: providerFactory,
		geminiEmbedding: geminiEmbedding,
		qdrantSearch:    qdrantSearch,
		faker:           faker,
		postgresRepo:    postgresRepo,
	}
}

// Health returns service health status
func (h *LLMHandler) Health(c *gin.Context) {
	provider := h.providerFactory.GetDefaultProviderName()
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "llm",
		"provider":  provider,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// ParseRequest parses a natural language request
func (h *LLMHandler) ParseRequest(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	requestIDStr, _ := requestID.(string)

	var req struct {
		NaturalLanguage string `json:"natural_language" binding:"required"`
		Provider        string `json:"provider,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithRequestID(requestIDStr).Debug().
			Err(err).
			Msg("Invalid parse request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "natural_language is required"})
		return
	}

	// Get LLM provider
	provider := h.providerFactory.GetProvider(req.Provider)
	if provider == nil {
		logger.WithRequestID(requestIDStr).Warn().
			Str("provider", req.Provider).
			Msg("No LLM provider available")
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No LLM provider available"})
		return
	}

	logger.WithRequestID(requestIDStr).Info().
		Str("provider", req.Provider).
		Int("nl_length", len(req.NaturalLanguage)).
		Msg("Parsing natural language request")

	// Get API context from vector search (RAG)
	apiContext, err := h.retrieveAPIContext(c.Request.Context(), req.NaturalLanguage)
	if err != nil {
		logger.WithRequestID(requestIDStr).Err(err).
			Msg("Failed to retrieve API context")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve API context"})
		return
	}

	// Build prompt and call LLM
	prompt := prompts.ParseRequestPrompt(req.NaturalLanguage, apiContext)
	response, err := provider.Complete(c.Request.Context(), prompt)
	if err != nil {
		logger.WithRequestID(requestIDStr).Err(err).
			Str("provider", req.Provider).
			Msg("LLM completion failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LLM error: " + err.Error()})
		return
	}

	logger.WithRequestID(requestIDStr).Debug().
		Str("provider", req.Provider).
		Int("response_length", len(response)).
		Msg("LLM response received")

	// Parse LLM response - try to extract JSON from markdown if needed
	var parseResult entities.ParseResult
	jsonStr := extractJSON(response)
	if err := json.Unmarshal([]byte(jsonStr), &parseResult); err != nil {
		// Fallback if JSON parsing fails
		parseResult = entities.ParseResult{
			Intent:     "Unable to parse LLM response",
			Confidence: 0.5,
		}
	}

	// Auto-generate test data for common fields that shouldn't require user input
	autoGenerateFields := map[string]bool{
		"card_number":    true,
		"card_expiry":    true,
		"expiry_month":   true,
		"expiry_year":    true,
		"cvv":            true,
		"cvc":            true,
		"pan":            true,
		"account_number": true,
		"first_name":     true,
		"last_name":      true,
		"address":        true,
		"city":           true,
		"country":        true,
		"postal_code":    true,
		"zip_code":       true,
		"phone":          true,
		"email":          true,
	}

	// Auto-fill parameters marked with "[AUTO]" or nil for auto-generatable fields
	if len(parseResult.Parameters) > 0 {
		for key, val := range parseResult.Parameters {
			shouldAutoGenerate := false
			
			// Check if LLM marked it as "[AUTO]"
			if strVal, ok := val.(string); ok && strings.ToUpper(strVal) == "[AUTO]" {
				shouldAutoGenerate = true
			}
			
			// Check if nil and it's an auto-generatable field
			if val == nil {
				keyLower := strings.ToLower(key)
				if autoGenerateFields[keyLower] {
					shouldAutoGenerate = true
				} else {
					// Also check for partial matches
					for autoField := range autoGenerateFields {
						if strings.Contains(keyLower, autoField) || strings.Contains(autoField, keyLower) {
							shouldAutoGenerate = true
							break
						}
					}
				}
			}
			
			if shouldAutoGenerate {
				// Auto-generate the value using faker
				parseResult.Parameters[key] = h.faker.GenerateByType(key, "string", "")
			}
		}
	}

	// Check if clarification is still needed (only for fields that couldn't be auto-generated)
	if len(parseResult.Parameters) > 0 {
		for key, val := range parseResult.Parameters {
			if val == nil {
				parseResult.NeedsClarify = true
				parseResult.Clarification = &entities.Clarification{
					ID:        uuid.New(),
					Message:   "Please provide value for: " + key,
					Type:      "free_text",
					FieldName: key,
				}
				break
			}
		}
	}

	c.JSON(http.StatusOK, parseResult)
}

// ConstructRequest constructs an API request from parse result
func (h *LLMHandler) ConstructRequest(c *gin.Context) {
	requestID, _ := c.Get("request_id")
	requestIDStr, _ := requestID.(string)

	var req struct {
		ParseResult  map[string]interface{} `json:"parse_result" binding:"required"`
		APIConfig    map[string]interface{} `json:"api_config,omitempty"`
		GenerateData bool                   `json:"generate_data"`
		Provider     string                 `json:"provider,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parse_result is required"})
		return
	}

	// Generate test data if requested
	generatedData := make(map[string]interface{})
	if req.GenerateData {
		if params, ok := req.ParseResult["parameters"].(map[string]interface{}); ok {
			for key, val := range params {
				if val == nil {
					generatedData[key] = h.faker.GenerateByType(key, "string", "")
				}
			}
		}
	}

	// Get LLM provider
	provider := h.providerFactory.GetProvider(req.Provider)
	if provider == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "No LLM provider available"})
		return
	}

	// Retrieve API context using RAG for the api_name from parse result
	var apiContext string
	if apiName, ok := req.ParseResult["api_name"].(string); ok && apiName != "" {
		apiContext, _ = h.retrieveAPIContext(c.Request.Context(), apiName)
		logger.WithRequestID(requestIDStr).Debug().
			Str("api_name", apiName).
			Str("api_context_length", fmt.Sprintf("%d", len(apiContext))).
			Msg("Retrieved API context for construct")
	}

	// Build prompt with API context included
	parseResultJSON, _ := json.Marshal(req.ParseResult)
	apiConfigJSON, _ := json.Marshal(req.APIConfig)
	prompt := prompts.ConstructRequestPromptWithContext(string(parseResultJSON), string(apiConfigJSON), apiContext, generatedData)

	// Call LLM
	response, err := provider.Complete(c.Request.Context(), prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LLM error: " + err.Error()})
		return
	}

	// Parse response into APICall - extract JSON from markdown if needed
	var apiCall entities.APICall
	jsonStr := extractJSON(response)
	
	// Check if we extracted valid JSON
	if jsonStr == "" || jsonStr == response {
		logger.WithRequestID(requestIDStr).Warn().
			Str("response_preview", truncateString(response, 200)).
			Msg("No valid JSON found in LLM response")
		c.JSON(http.StatusOK, gin.H{
			"api_call":       nil,
			"raw_json":       response,
			"parse_error":    "No valid JSON structure found in LLM response",
			"generated_data": generatedData,
		})
		return
	}

	if err := json.Unmarshal([]byte(jsonStr), &apiCall); err != nil {
		// Return partial success with raw JSON for debugging
		logger.WithRequestID(requestIDStr).Warn().
			Err(err).
			Str("json_preview", truncateString(jsonStr, 200)).
			Msg("Failed to unmarshal LLM JSON response")
		c.JSON(http.StatusOK, gin.H{
			"api_call":       nil,
			"raw_json":       jsonStr,
			"parse_error":    err.Error(),
			"generated_data": generatedData,
		})
		return
	}

	apiCall.ID = uuid.New()

	c.JSON(http.StatusOK, gin.H{
		"api_call":       apiCall,
		"generated_data": generatedData,
	})
}

// Clarify handles clarification responses
func (h *LLMHandler) Clarify(c *gin.Context) {
	var req entities.ClarificationResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Return the value to be used
	value := req.SelectedValue
	if value == "" {
		value = req.FreeText
	}

	c.JSON(http.StatusOK, gin.H{
		"clarification_id": req.ClarificationID,
		"value":            value,
		"status":           "resolved",
	})
}

// GenerateData generates test data
func (h *LLMHandler) GenerateData(c *gin.Context) {
	var req entities.GenerateDataRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	value := h.faker.GenerateByType(req.FieldName, req.FieldType, req.Format)

	c.JSON(http.StatusOK, gin.H{
		"field_name": req.FieldName,
		"field_type": req.FieldType,
		"value":      value,
	})
}

// Learn records a successful test pattern
func (h *LLMHandler) Learn(c *gin.Context) {
	var req struct {
		APISpecID   uuid.UUID              `json:"api_spec_id" binding:"required"`
		PatternData map[string]interface{} `json:"pattern_data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.postgresRepo.IncrementSuccessCount(c.Request.Context(), req.APISpecID, req.PatternData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record pattern"})
		return
	}

	// Check if we've learned enough
	learned, _ := h.postgresRepo.CheckIfLearnedEnough(c.Request.Context(), req.APISpecID)

	c.JSON(http.StatusOK, gin.H{
		"status":          "recorded",
		"pattern_learned": learned,
	})
}

// ListProviders lists available LLM providers
func (h *LLMHandler) ListProviders(c *gin.Context) {
	providers := h.providerFactory.ListAvailableProviders()
	defaultProvider := h.providerFactory.GetDefaultProviderName()

	c.JSON(http.StatusOK, gin.H{
		"providers":        providers,
		"default_provider": defaultProvider,
	})
}

// retrieveAPIContext retrieves relevant API context using RAG
func (h *LLMHandler) retrieveAPIContext(ctx context.Context, query string) (string, error) {
	// Generate embedding for query using Gemini
	if h.geminiEmbedding == nil || !h.geminiEmbedding.IsAvailable() {
		return "No API context available (embeddings not configured)", nil
	}

	embedding, err := h.geminiEmbedding.GenerateEmbedding(ctx, query)
	if err != nil {
		return "No API context available (embedding failed: " + err.Error() + ")", nil
	}

	// Search Qdrant
	results, err := h.qdrantSearch.Search(embedding, 3)
	if err != nil {
		return "No API context available (search failed: " + err.Error() + ")", nil
	}

	// Build context from results
	var contexts []map[string]interface{}
	for _, r := range results {
		contexts = append(contexts, map[string]interface{}{
			"api_name":    r.APIName,
			"version":     r.Version,
			"description": r.Description,
			"config":      r.Config,
		})
	}

	return prompts.BuildAPIContext(contexts), nil
}

// extractJSON attempts to extract JSON from a response that may be wrapped in markdown
func extractJSON(response string) string {
	response = strings.TrimSpace(response)

	// Try to extract JSON from markdown code block
	if strings.Contains(response, "```json") {
		start := strings.Index(response, "```json") + 7
		end := strings.LastIndex(response, "```")
		if start < end {
			return strings.TrimSpace(response[start:end])
		}
	}

	// Try to extract JSON from generic code block
	if strings.Contains(response, "```") {
		start := strings.Index(response, "```") + 3
		// Skip language identifier if present
		if idx := strings.Index(response[start:], "\n"); idx != -1 {
			start += idx + 1
		}
		end := strings.LastIndex(response, "```")
		if start < end {
			return strings.TrimSpace(response[start:end])
		}
	}

	// Try to find JSON object directly
	if start := strings.Index(response, "{"); start != -1 {
		if end := strings.LastIndex(response, "}"); end > start {
			return response[start : end+1]
		}
	}

	return response
}

// truncateString truncates a string to a maximum length for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
