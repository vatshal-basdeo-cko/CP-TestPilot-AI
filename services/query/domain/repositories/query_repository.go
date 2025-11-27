package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/testpilot-ai/query/domain/entities"
)

// QueryRepository defines the interface for querying test data
type QueryRepository interface {
	// FindExecutionByID retrieves a single execution
	FindExecutionByID(ctx context.Context, id uuid.UUID) (*entities.TestExecution, error)
	
	// ListExecutions retrieves executions with filters
	ListExecutions(ctx context.Context, filters Filters) ([]entities.TestExecution, int64, error)
	
	// GetAnalytics retrieves aggregated statistics
	GetAnalytics(ctx context.Context, startDate, endDate *time.Time) (*entities.Analytics, error)
	
	// GetAPIAnalytics retrieves per-API statistics
	GetAPIAnalytics(ctx context.Context, apiSpecID uuid.UUID) (*entities.APIStats, error)
}

// Filters for querying executions
type Filters struct {
	UserID    *uuid.UUID
	APISpecID *uuid.UUID
	Status    string
	StartDate *time.Time
	EndDate   *time.Time
	Search    string
	Limit     int
	Offset    int
}

