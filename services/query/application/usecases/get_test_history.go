package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/testpilot-ai/query/domain/entities"
	"github.com/testpilot-ai/query/domain/repositories"
)

// GetTestHistoryUseCase handles test history retrieval
type GetTestHistoryUseCase struct {
	repo repositories.QueryRepository
}

// NewGetTestHistoryUseCase creates a new use case
func NewGetTestHistoryUseCase(repo repositories.QueryRepository) *GetTestHistoryUseCase {
	return &GetTestHistoryUseCase{repo: repo}
}

// Execute retrieves test execution history
func (uc *GetTestHistoryUseCase) Execute(ctx context.Context, filters repositories.Filters) ([]entities.TestExecution, int64, error) {
	return uc.repo.ListExecutions(ctx, filters)
}

// GetByID retrieves a single execution
func (uc *GetTestHistoryUseCase) GetByID(ctx context.Context, id uuid.UUID) (*entities.TestExecution, error) {
	return uc.repo.FindExecutionByID(ctx, id)
}

