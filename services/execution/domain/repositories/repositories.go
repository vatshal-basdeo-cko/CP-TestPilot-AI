package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/testpilot-ai/execution/domain/entities"
)

// ExecutionRepository defines the interface for execution data operations
type ExecutionRepository interface {
	// SaveExecution saves an execution record
	SaveExecution(ctx context.Context, request *entities.APIRequest, response *entities.APIResponse) error
	
	// FindExecutionByID retrieves an execution by ID
	FindExecutionByID(ctx context.Context, id uuid.UUID) (*entities.APIResponse, error)
	
	// ListExecutions retrieves executions with pagination
	ListExecutions(ctx context.Context, limit, offset int) ([]*entities.APIResponse, error)
}

// EnvironmentRepository defines the interface for environment operations
type EnvironmentRepository interface {
	// CreateEnvironment creates a new environment
	CreateEnvironment(ctx context.Context, env *entities.Environment) error
	
	// FindEnvironmentByID retrieves an environment by ID
	FindEnvironmentByID(ctx context.Context, id uuid.UUID) (*entities.Environment, error)
	
	// FindEnvironmentByName retrieves an environment by name
	FindEnvironmentByName(ctx context.Context, name string) (*entities.Environment, error)
	
	// UpdateEnvironment updates an environment
	UpdateEnvironment(ctx context.Context, env *entities.Environment) error
	
	// DeleteEnvironment deletes an environment
	DeleteEnvironment(ctx context.Context, id uuid.UUID) error
	
	// ListEnvironments retrieves all environments
	ListEnvironments(ctx context.Context) ([]*entities.Environment, error)
}

