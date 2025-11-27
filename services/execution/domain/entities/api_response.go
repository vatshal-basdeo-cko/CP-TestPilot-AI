package entities

import (
	"time"

	"github.com/google/uuid"
)

// APIResponse represents the response from an API execution
type APIResponse struct {
	ID              uuid.UUID              `json:"id"`
	RequestID       uuid.UUID              `json:"request_id"`
	StatusCode      int                    `json:"status_code"`
	Headers         map[string][]string    `json:"headers"`
	Body            interface{}            `json:"body"`
	ExecutionTimeMs int64                  `json:"execution_time_ms"`
	Error           string                 `json:"error,omitempty"`
	Success         bool                   `json:"success"`
	Timestamp       time.Time              `json:"timestamp"`
}

// NewAPIResponse creates a new API response entity
func NewAPIResponse(requestID uuid.UUID) *APIResponse {
	return &APIResponse{
		ID:        uuid.New(),
		RequestID: requestID,
		Headers:   make(map[string][]string),
		Timestamp: time.Now(),
	}
}

// IsSuccessful checks if the response indicates success
func (r *APIResponse) IsSuccessful() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

