package entities

import (
	"time"

	"github.com/google/uuid"
)

// APIRequest represents an API request to be executed
type APIRequest struct {
	ID                     uuid.UUID              `json:"id"`
	Method                 string                 `json:"method"`
	URL                    string                 `json:"url"`
	Headers                map[string]string      `json:"headers"`
	QueryParams            map[string]interface{} `json:"query_params"`
	Body                   interface{}            `json:"body,omitempty"`
	Timeout                int                    `json:"timeout"` // in seconds
	APISpecID              *uuid.UUID             `json:"api_spec_id,omitempty"`
	APIName                string                 `json:"api_name,omitempty"`
	EndpointName           string                 `json:"endpoint_name,omitempty"`
	UserID                 *uuid.UUID             `json:"user_id,omitempty"`
	NaturalLanguageRequest string                 `json:"natural_language_request,omitempty"`
	CreatedAt              time.Time              `json:"created_at"`
}

// NewAPIRequest creates a new API request entity
func NewAPIRequest(method, url string) *APIRequest {
	return &APIRequest{
		ID:          uuid.New(),
		Method:      method,
		URL:         url,
		Headers:     make(map[string]string),
		QueryParams: make(map[string]interface{}),
		Timeout:     30, // default 30 seconds
		CreatedAt:   time.Now(),
	}
}

// Validate checks if the request is valid
func (r *APIRequest) Validate() error {
	if r.Method == "" {
		return ErrInvalidMethod
	}
	if r.URL == "" {
		return ErrInvalidURL
	}
	validMethods := []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	valid := false
	for _, m := range validMethods {
		if r.Method == m {
			valid = true
			break
		}
	}
	if !valid {
		return ErrInvalidMethod
	}
	return nil
}

