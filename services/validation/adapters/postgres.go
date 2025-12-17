package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testpilot-ai/validation/domain/entities"
)

// PostgresRepository handles database operations for validation
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// GetRulesForAPI retrieves validation rules for a specific API
func (r *PostgresRepository) GetRulesForAPI(ctx context.Context, apiSpecID uuid.UUID) ([]entities.ValidationRule, error) {
	query := `
		SELECT id, api_spec_id, rule_type, rule_definition, created_at, updated_at
		FROM validation_rules
		WHERE api_spec_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.pool.Query(ctx, query, apiSpecID)
	if err != nil {
		return nil, fmt.Errorf("failed to query rules: %w", err)
	}
	defer rows.Close()

	var rules []entities.ValidationRule
	for rows.Next() {
		var rule entities.ValidationRule
		err := rows.Scan(
			&rule.ID,
			&rule.APISpecID,
			&rule.RuleType,
			&rule.RuleDefinition,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// GetAllRules retrieves all validation rules
func (r *PostgresRepository) GetAllRules(ctx context.Context) ([]entities.ValidationRule, error) {
	query := `
		SELECT id, api_spec_id, rule_type, rule_definition, created_at, updated_at
		FROM validation_rules
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query rules: %w", err)
	}
	defer rows.Close()

	var rules []entities.ValidationRule
	for rows.Next() {
		var rule entities.ValidationRule
		err := rows.Scan(
			&rule.ID,
			&rule.APISpecID,
			&rule.RuleType,
			&rule.RuleDefinition,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		rules = append(rules, rule)
	}

	return rules, nil
}

// CreateRule creates a new validation rule
func (r *PostgresRepository) CreateRule(ctx context.Context, rule *entities.ValidationRule) error {
	query := `
		INSERT INTO validation_rules (id, api_spec_id, rule_type, rule_definition, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	now := time.Now()
	rule.ID = uuid.New()
	rule.CreatedAt = now
	rule.UpdatedAt = now

	_, err := r.pool.Exec(ctx, query,
		rule.ID,
		rule.APISpecID,
		rule.RuleType,
		rule.RuleDefinition,
		rule.CreatedAt,
		rule.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create rule: %w", err)
	}

	return nil
}

// UpdateRule updates an existing validation rule
func (r *PostgresRepository) UpdateRule(ctx context.Context, rule *entities.ValidationRule) error {
	query := `
		UPDATE validation_rules
		SET rule_type = $2, rule_definition = $3, updated_at = $4
		WHERE id = $1
	`

	rule.UpdatedAt = time.Now()

	_, err := r.pool.Exec(ctx, query,
		rule.ID,
		rule.RuleType,
		rule.RuleDefinition,
		rule.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update rule: %w", err)
	}

	return nil
}

// DeleteRule deletes a validation rule
func (r *PostgresRepository) DeleteRule(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM validation_rules WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}
	return nil
}

