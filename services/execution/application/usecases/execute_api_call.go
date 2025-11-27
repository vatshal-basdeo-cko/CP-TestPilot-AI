package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/testpilot-ai/execution/domain/entities"
	"github.com/testpilot-ai/execution/domain/repositories"
)

// ExecuteAPICallUseCase handles API execution logic
type ExecuteAPICallUseCase struct {
	executionRepo repositories.ExecutionRepository
	httpClient    *http.Client
}

// NewExecuteAPICallUseCase creates a new use case instance
func NewExecuteAPICallUseCase(repo repositories.ExecutionRepository) *ExecuteAPICallUseCase {
	return &ExecuteAPICallUseCase{
		executionRepo: repo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Execute executes an API call
func (uc *ExecuteAPICallUseCase) Execute(ctx context.Context, request *entities.APIRequest) (*entities.APIResponse, error) {
	// Validate request
	if err := request.Validate(); err != nil {
		return nil, err
	}

	// Prepare response
	response := entities.NewAPIResponse(request.ID)
	startTime := time.Now()

	// Execute HTTP request
	httpReq, err := uc.buildHTTPRequest(ctx, request)
	if err != nil {
		response.Error = err.Error()
		response.Success = false
		return response, err
	}

	// Set timeout
	if request.Timeout > 0 {
		uc.httpClient.Timeout = time.Duration(request.Timeout) * time.Second
	}

	// Make the request
	httpResp, err := uc.httpClient.Do(httpReq)
	if err != nil {
		response.Error = err.Error()
		response.Success = false
		response.ExecutionTimeMs = time.Since(startTime).Milliseconds()
		
		// Save failed execution
		_ = uc.executionRepo.SaveExecution(ctx, request, response)
		return response, err
	}
	defer httpResp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		response.Error = fmt.Sprintf("failed to read response body: %v", err)
		response.Success = false
		response.ExecutionTimeMs = time.Since(startTime).Milliseconds()
		return response, err
	}

	// Parse response
	response.StatusCode = httpResp.StatusCode
	response.Headers = httpResp.Header
	response.ExecutionTimeMs = time.Since(startTime).Milliseconds()

	// Try to parse as JSON
	var bodyJSON interface{}
	if err := json.Unmarshal(bodyBytes, &bodyJSON); err == nil {
		response.Body = bodyJSON
	} else {
		response.Body = string(bodyBytes)
	}

	response.Success = response.IsSuccessful()

	// Save execution to database
	if err := uc.executionRepo.SaveExecution(ctx, request, response); err != nil {
		// Log error but don't fail the execution
		fmt.Printf("Failed to save execution: %v\n", err)
	}

	return response, nil
}

// buildHTTPRequest creates an HTTP request from APIRequest entity
func (uc *ExecuteAPICallUseCase) buildHTTPRequest(ctx context.Context, request *entities.APIRequest) (*http.Request, error) {
	// Build URL with query params
	url := request.URL
	if len(request.QueryParams) > 0 {
		params := []string{}
		for k, v := range request.QueryParams {
			params = append(params, fmt.Sprintf("%s=%v", k, v))
		}
		if strings.Contains(url, "?") {
			url = url + "&" + strings.Join(params, "&")
		} else {
			url = url + "?" + strings.Join(params, "&")
		}
	}

	// Prepare body
	var bodyReader io.Reader
	if request.Body != nil {
		bodyBytes, err := json.Marshal(request.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = strings.NewReader(string(bodyBytes))
	}

	// Create request
	httpReq, err := http.NewRequestWithContext(ctx, request.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for k, v := range request.Headers {
		httpReq.Header.Set(k, v)
	}

	// Default content type if not set
	if request.Body != nil && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	return httpReq, nil
}

