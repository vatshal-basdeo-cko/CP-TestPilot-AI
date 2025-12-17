package entities

import (
	"time"

	"github.com/google/uuid"
)

// TestRequest represents a natural language test request
type TestRequest struct {
	ID              uuid.UUID `json:"id"`
	NaturalLanguage string    `json:"natural_language"`
	UserID          uuid.UUID `json:"user_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// APICall represents a constructed API call
type APICall struct {
	ID              uuid.UUID              `json:"id"`
	Method          string                 `json:"method"`
	URL             string                 `json:"url"`
	Path            string                 `json:"path"`
	Headers         map[string]string      `json:"headers,omitempty"`
	QueryParams     map[string]string      `json:"query_params,omitempty"`
	Body            map[string]interface{} `json:"body,omitempty"`
	APISpecID       uuid.UUID              `json:"api_spec_id,omitempty"`
	APIName         string                 `json:"api_name,omitempty"`
	EndpointName    string                 `json:"endpoint_name,omitempty"`
	Confidence      float64                `json:"confidence"`
}

// RetrievalContext represents context retrieved from vector search
type RetrievalContext struct {
	APIName     string                 `json:"api_name"`
	Version     string                 `json:"version"`
	Description string                 `json:"description"`
	Endpoints   []EndpointContext      `json:"endpoints"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Score       float32                `json:"score"`
}

// EndpointContext represents a single endpoint's context
type EndpointContext struct {
	Name        string                 `json:"name"`
	Path        string                 `json:"path"`
	Method      string                 `json:"method"`
	Description string                 `json:"description"`
	Parameters  []ParameterContext     `json:"parameters,omitempty"`
	Schema      map[string]interface{} `json:"schema,omitempty"`
}

// ParameterContext represents a parameter's context
type ParameterContext struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
}

// Clarification represents a clarification request
type Clarification struct {
	ID        uuid.UUID           `json:"id"`
	Message   string              `json:"message"`
	Type      string              `json:"type"` // multiple_choice, free_text
	Options   []ClarificationOption `json:"options,omitempty"`
	FieldName string              `json:"field_name,omitempty"`
}

// ClarificationOption represents an option for clarification
type ClarificationOption struct {
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

// ClarificationResponse represents user's response to clarification
type ClarificationResponse struct {
	ClarificationID uuid.UUID `json:"clarification_id"`
	SelectedValue   string    `json:"selected_value,omitempty"`
	FreeText        string    `json:"free_text,omitempty"`
}

// ParseResult represents the result of parsing natural language
type ParseResult struct {
	Intent       string                 `json:"intent"`
	APIName      string                 `json:"api_name,omitempty"`
	Endpoint     string                 `json:"endpoint,omitempty"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
	Confidence   float64                `json:"confidence"`
	NeedsClarify bool                   `json:"needs_clarification"`
	Clarification *Clarification        `json:"clarification,omitempty"`
}

// LearnedPattern represents a learned successful test pattern
type LearnedPattern struct {
	ID           uuid.UUID              `json:"id"`
	APISpecID    uuid.UUID              `json:"api_spec_id"`
	PatternData  map[string]interface{} `json:"pattern_data"`
	SuccessCount int                    `json:"success_count"`
	LastUpdated  time.Time              `json:"last_updated"`
}

// GenerateDataRequest represents a request to generate test data
type GenerateDataRequest struct {
	FieldName string `json:"field_name"`
	FieldType string `json:"field_type"`
	Format    string `json:"format,omitempty"`
	Locale    string `json:"locale,omitempty"`
}

