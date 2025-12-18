package entities

import (
	"time"

	"github.com/google/uuid"
)

// APISpecification represents an ingested API configuration
type APISpecification struct {
	ID          uuid.UUID              `json:"id"`
	Name        string                 `json:"name"`
	Version     string                 `json:"version"`
	SourceType  string                 `json:"source_type"`
	SourcePath  string                 `json:"source_path,omitempty"`
	ContentHash string                 `json:"content_hash"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CreatedBy   *uuid.UUID             `json:"created_by,omitempty"`
}

// IngestionResult represents the result of an ingestion operation
type IngestionResult struct {
	ID           uuid.UUID `json:"id"`
	SourceType   string    `json:"source_type"`
	SourcePath   string    `json:"source_path,omitempty"`
	Status       string    `json:"status"`
	APIsIngested int       `json:"apis_ingested"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

// APIEndpoint represents a single API endpoint from the config
type APIEndpoint struct {
	Name                string                   `yaml:"name" json:"name"`
	Path                string                   `yaml:"path" json:"path"`
	Method              string                   `yaml:"method" json:"method"`
	Description         string                   `yaml:"description" json:"description"`
	Authentication      *AuthConfig              `yaml:"authentication" json:"authentication,omitempty"`
	Parameters          []Parameter              `yaml:"parameters" json:"parameters,omitempty"`
	RequestSchema       map[string]interface{}   `yaml:"request_schema" json:"request_schema,omitempty"`
	ResponseSchema      map[string]interface{}   `yaml:"response_schema" json:"response_schema,omitempty"`
	ExpectedStatusCodes []map[string]interface{} `yaml:"expected_status_codes" json:"expected_status_codes,omitempty"`
	Examples            []Example                `yaml:"examples" json:"examples,omitempty"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type   string `yaml:"type" json:"type"`
	Header string `yaml:"header" json:"header,omitempty"`
	Key    string `yaml:"key" json:"key,omitempty"`
}

// Parameter represents an API parameter
type Parameter struct {
	Name        string `yaml:"name" json:"name"`
	Type        string `yaml:"type" json:"type"`
	In          string `yaml:"in" json:"in,omitempty"`
	Required    bool   `yaml:"required" json:"required"`
	Description string `yaml:"description" json:"description,omitempty"`
	Default     string `yaml:"default" json:"default,omitempty"`
	Format      string `yaml:"format" json:"format,omitempty"`
	Example     string `yaml:"example" json:"example,omitempty"`
}

// Example represents a request/response example
type Example struct {
	Name     string                 `yaml:"name" json:"name"`
	Request  map[string]interface{} `yaml:"request" json:"request"`
	Response map[string]interface{} `yaml:"response" json:"response"`
}

// APIConfig represents a full API configuration file
type APIConfig struct {
	Name        string        `yaml:"name" json:"name"`
	Version     string        `yaml:"version" json:"version"`
	Description string        `yaml:"description" json:"description"`
	BaseURL     string        `yaml:"base_url" json:"base_url"`
	Endpoints   []APIEndpoint `yaml:"endpoints" json:"endpoints"`
}

// PostmanCollection represents a Postman collection
type PostmanCollection struct {
	Info struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Schema      string `json:"schema"`
	} `json:"info"`
	Item []PostmanItem `json:"item"`
}

// PostmanItem represents an item in a Postman collection
type PostmanItem struct {
	Name    string          `json:"name"`
	Request *PostmanRequest `json:"request,omitempty"`
	Item    []PostmanItem   `json:"item,omitempty"` // For folders
}

// PostmanRequest represents a Postman request
type PostmanRequest struct {
	Method string `json:"method"`
	Header []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"header"`
	URL  PostmanURL   `json:"url"`
	Body *PostmanBody `json:"body,omitempty"`
}

// PostmanURL represents a Postman URL
type PostmanURL struct {
	Raw   string   `json:"raw"`
	Host  []string `json:"host"`
	Path  []string `json:"path"`
	Query []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"query"`
}

// PostmanBody represents a Postman request body
type PostmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}
