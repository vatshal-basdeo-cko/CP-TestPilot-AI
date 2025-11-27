package adapters

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testpilot-ai/execution/domain/entities"
)

// EnvironmentRepository implements environment repository using PostgreSQL
type EnvironmentRepository struct {
	pool *pgxpool.Pool
}

// NewEnvironmentRepository creates a new environment repository
func NewEnvironmentRepository(pool *pgxpool.Pool) *EnvironmentRepository {
	return &EnvironmentRepository{
		pool: pool,
	}
}

// CreateEnvironment creates a new environment
func (r *EnvironmentRepository) CreateEnvironment(ctx context.Context, env *entities.Environment) error {
	query := `
		INSERT INTO environments (id, name, base_url, auth_config, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	authConfig, err := json.Marshal(env.AuthConfig)
	if err != nil {
		return err
	}

	_, err = r.pool.Exec(ctx, query,
		env.ID,
		env.Name,
		env.BaseURL,
		authConfig,
		env.CreatedAt,
		env.UpdatedAt,
	)

	return err
}

// FindEnvironmentByID retrieves an environment by ID
func (r *EnvironmentRepository) FindEnvironmentByID(ctx context.Context, id uuid.UUID) (*entities.Environment, error) {
	query := `
		SELECT id, name, base_url, auth_config, created_at, updated_at
		FROM environments
		WHERE id = $1
	`

	var env entities.Environment
	var authConfigJSON []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&env.ID,
		&env.Name,
		&env.BaseURL,
		&authConfigJSON,
		&env.CreatedAt,
		&env.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(authConfigJSON, &env.AuthConfig); err != nil {
		env.AuthConfig = make(map[string]interface{})
	}

	env.Active = true
	return &env, nil
}

// FindEnvironmentByName retrieves an environment by name
func (r *EnvironmentRepository) FindEnvironmentByName(ctx context.Context, name string) (*entities.Environment, error) {
	query := `
		SELECT id, name, base_url, auth_config, created_at, updated_at
		FROM environments
		WHERE name = $1
	`

	var env entities.Environment
	var authConfigJSON []byte

	err := r.pool.QueryRow(ctx, query, name).Scan(
		&env.ID,
		&env.Name,
		&env.BaseURL,
		&authConfigJSON,
		&env.CreatedAt,
		&env.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(authConfigJSON, &env.AuthConfig); err != nil {
		env.AuthConfig = make(map[string]interface{})
	}

	env.Active = true
	return &env, nil
}

// UpdateEnvironment updates an environment
func (r *EnvironmentRepository) UpdateEnvironment(ctx context.Context, env *entities.Environment) error {
	query := `
		UPDATE environments
		SET name = $2, base_url = $3, auth_config = $4, updated_at = $5
		WHERE id = $1
	`

	authConfig, err := json.Marshal(env.AuthConfig)
	if err != nil {
		return err
	}

	env.UpdatedAt = time.Now()

	_, err = r.pool.Exec(ctx, query,
		env.ID,
		env.Name,
		env.BaseURL,
		authConfig,
		env.UpdatedAt,
	)

	return err
}

// DeleteEnvironment deletes an environment
func (r *EnvironmentRepository) DeleteEnvironment(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM environments WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	return err
}

// ListEnvironments retrieves all environments
func (r *EnvironmentRepository) ListEnvironments(ctx context.Context) ([]*entities.Environment, error) {
	query := `
		SELECT id, name, base_url, auth_config, created_at, updated_at
		FROM environments
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var environments []*entities.Environment
	for rows.Next() {
		var env entities.Environment
		var authConfigJSON []byte

		err := rows.Scan(
			&env.ID,
			&env.Name,
			&env.BaseURL,
			&authConfigJSON,
			&env.CreatedAt,
			&env.UpdatedAt,
		)
		if err != nil {
			continue
		}

		if err := json.Unmarshal(authConfigJSON, &env.AuthConfig); err != nil {
			env.AuthConfig = make(map[string]interface{})
		}

		env.Active = true
		environments = append(environments, &env)
	}

	return environments, nil
}

