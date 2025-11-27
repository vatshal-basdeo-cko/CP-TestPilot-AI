package entities

import (
	"time"

	"github.com/google/uuid"
)

// TestExecution represents a test execution record
type TestExecution struct {
	ID                      uuid.UUID              `json:"id"`
	UserID                  *uuid.UUID             `json:"user_id"`
	APISpecID               *uuid.UUID             `json:"api_spec_id"`
	NaturalLanguageRequest  string                 `json:"natural_language_request"`
	ConstructedRequest      map[string]interface{} `json:"constructed_request"`
	Response                map[string]interface{} `json:"response"`
	ValidationResult        map[string]interface{} `json:"validation_result"`
	Status                  string                 `json:"status"` // success, failed, error
	ExecutionTimeMs         int64                  `json:"execution_time_ms"`
	CreatedAt               time.Time              `json:"created_at"`
}

// Analytics represents aggregated statistics
type Analytics struct {
	TotalTests       int64              `json:"total_tests"`
	SuccessfulTests  int64              `json:"successful_tests"`
	FailedTests      int64              `json:"failed_tests"`
	SuccessRate      float64            `json:"success_rate"`
	AvgExecutionTime float64            `json:"avg_execution_time_ms"`
	TopAPIs          []APIStats         `json:"top_apis"`
	RecentTests      []TestExecution    `json:"recent_tests"`
}

// APIStats represents statistics for a specific API
type APIStats struct {
	APIName      string  `json:"api_name"`
	TestCount    int64   `json:"test_count"`
	SuccessCount int64   `json:"success_count"`
	SuccessRate  float64 `json:"success_rate"`
}

