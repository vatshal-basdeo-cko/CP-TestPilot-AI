package usecases

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/testpilot-ai/query/domain/entities"
	"github.com/testpilot-ai/query/domain/repositories"
)

// GetAnalyticsUseCase handles analytics retrieval
type GetAnalyticsUseCase struct {
	repo repositories.QueryRepository
}

// NewGetAnalyticsUseCase creates a new use case
func NewGetAnalyticsUseCase(repo repositories.QueryRepository) *GetAnalyticsUseCase {
	return &GetAnalyticsUseCase{repo: repo}
}

// GetOverview retrieves overall analytics
func (uc *GetAnalyticsUseCase) GetOverview(ctx context.Context, startDate, endDate *time.Time) (*entities.Analytics, error) {
	return uc.repo.GetAnalytics(ctx, startDate, endDate)
}

// GetByAPI retrieves analytics for specific API
func (uc *GetAnalyticsUseCase) GetByAPI(ctx context.Context, apiSpecID uuid.UUID) (*entities.APIStats, error) {
	return uc.repo.GetAPIAnalytics(ctx, apiSpecID)
}




