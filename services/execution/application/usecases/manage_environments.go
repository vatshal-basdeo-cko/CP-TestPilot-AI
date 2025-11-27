package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/testpilot-ai/execution/domain/entities"
	"github.com/testpilot-ai/execution/domain/repositories"
)

// ManageEnvironmentsUseCase handles environment management logic
type ManageEnvironmentsUseCase struct {
	envRepo repositories.EnvironmentRepository
}

// NewManageEnvironmentsUseCase creates a new use case instance
func NewManageEnvironmentsUseCase(repo repositories.EnvironmentRepository) *ManageEnvironmentsUseCase {
	return &ManageEnvironmentsUseCase{
		envRepo: repo,
	}
}

// CreateEnvironment creates a new environment
func (uc *ManageEnvironmentsUseCase) CreateEnvironment(ctx context.Context, env *entities.Environment) error {
	if err := env.Validate(); err != nil {
		return err
	}
	return uc.envRepo.CreateEnvironment(ctx, env)
}

// GetEnvironmentByID retrieves an environment by ID
func (uc *ManageEnvironmentsUseCase) GetEnvironmentByID(ctx context.Context, id uuid.UUID) (*entities.Environment, error) {
	return uc.envRepo.FindEnvironmentByID(ctx, id)
}

// GetEnvironmentByName retrieves an environment by name
func (uc *ManageEnvironmentsUseCase) GetEnvironmentByName(ctx context.Context, name string) (*entities.Environment, error) {
	return uc.envRepo.FindEnvironmentByName(ctx, name)
}

// UpdateEnvironment updates an environment
func (uc *ManageEnvironmentsUseCase) UpdateEnvironment(ctx context.Context, env *entities.Environment) error {
	if err := env.Validate(); err != nil {
		return err
	}
	return uc.envRepo.UpdateEnvironment(ctx, env)
}

// DeleteEnvironment deletes an environment
func (uc *ManageEnvironmentsUseCase) DeleteEnvironment(ctx context.Context, id uuid.UUID) error {
	return uc.envRepo.DeleteEnvironment(ctx, id)
}

// ListEnvironments retrieves all environments
func (uc *ManageEnvironmentsUseCase) ListEnvironments(ctx context.Context) ([]*entities.Environment, error) {
	return uc.envRepo.ListEnvironments(ctx)
}

