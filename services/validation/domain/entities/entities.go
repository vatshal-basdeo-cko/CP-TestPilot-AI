package entities

import (
	"time"

	"github.com/google/uuid"
)

// ValidationRule represents a validation rule for an API
type ValidationRule struct {
	ID             uuid.UUID              `json:"id"`
	APISpecID      uuid.UUID              `json:"api_spec_id"`
	RuleType       string                 `json:"rule_type"` // schema, status, custom
	RuleDefinition map[string]interface{} `json:"rule_definition"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// ValidationResult represents the result of a validation
type ValidationResult struct {
	IsValid       bool              `json:"is_valid"`
	StatusCheck   *StatusCheckResult `json:"status_check,omitempty"`
	SchemaCheck   *SchemaCheckResult `json:"schema_check,omitempty"`
	CustomChecks  []CustomCheckResult `json:"custom_checks,omitempty"`
	Errors        []string          `json:"errors,omitempty"`
	Warnings      []string          `json:"warnings,omitempty"`
	ValidatedAt   time.Time         `json:"validated_at"`
}

// StatusCheckResult represents status code validation result
type StatusCheckResult struct {
	Expected int  `json:"expected"`
	Actual   int  `json:"actual"`
	IsValid  bool `json:"is_valid"`
}

// SchemaCheckResult represents JSON schema validation result
type SchemaCheckResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors,omitempty"`
}

// CustomCheckResult represents a custom rule check result
type CustomCheckResult struct {
	RuleName string `json:"rule_name"`
	IsValid  bool   `json:"is_valid"`
	Message  string `json:"message,omitempty"`
}

// ValidationRequest represents a validation request
type ValidationRequest struct {
	APISpecID       *uuid.UUID             `json:"api_spec_id,omitempty"`
	Response        map[string]interface{} `json:"response"`
	StatusCode      int                    `json:"status_code"`
	ExpectedStatus  int                    `json:"expected_status,omitempty"`
	ExpectedSchema  map[string]interface{} `json:"expected_schema,omitempty"`
	PreviousSuccess map[string]interface{} `json:"previous_success,omitempty"` // For comparison
}

// DiffResult represents differences between two responses
type DiffResult struct {
	HasDifferences bool          `json:"has_differences"`
	Additions      []DiffEntry   `json:"additions,omitempty"`
	Deletions      []DiffEntry   `json:"deletions,omitempty"`
	Modifications  []DiffEntry   `json:"modifications,omitempty"`
}

// DiffEntry represents a single difference
type DiffEntry struct {
	Path     string      `json:"path"`
	OldValue interface{} `json:"old_value,omitempty"`
	NewValue interface{} `json:"new_value,omitempty"`
}

