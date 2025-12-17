package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testpilot-ai/llm/domain/entities"
)

// PostgresRepository handles database operations for LLM service
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// GetLearnedPatterns retrieves learned patterns for an API
func (r *PostgresRepository) GetLearnedPatterns(ctx context.Context, apiSpecID uuid.UUID) ([]entities.LearnedPattern, error) {
	query := `
		SELECT id, api_spec_id, pattern_data, success_count, last_updated
		FROM learned_patterns
		WHERE api_spec_id = $1
		ORDER BY success_count DESC
	`

	rows, err := r.pool.Query(ctx, query, apiSpecID)
	if err != nil {
		return nil, fmt.Errorf("failed to query patterns: %w", err)
	}
	defer rows.Close()

	var patterns []entities.LearnedPattern
	for rows.Next() {
		var pattern entities.LearnedPattern
		err := rows.Scan(
			&pattern.ID,
			&pattern.APISpecID,
			&pattern.PatternData,
			&pattern.SuccessCount,
			&pattern.LastUpdated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		patterns = append(patterns, pattern)
	}

	return patterns, nil
}

// SaveLearnedPattern saves or updates a learned pattern
func (r *PostgresRepository) SaveLearnedPattern(ctx context.Context, pattern *entities.LearnedPattern) error {
	query := `
		INSERT INTO learned_patterns (id, api_spec_id, pattern_data, success_count, last_updated)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			pattern_data = EXCLUDED.pattern_data,
			success_count = EXCLUDED.success_count,
			last_updated = EXCLUDED.last_updated
	`

	if pattern.ID == uuid.Nil {
		pattern.ID = uuid.New()
	}
	pattern.LastUpdated = time.Now()

	_, err := r.pool.Exec(ctx, query,
		pattern.ID,
		pattern.APISpecID,
		pattern.PatternData,
		pattern.SuccessCount,
		pattern.LastUpdated,
	)

	if err != nil {
		return fmt.Errorf("failed to save pattern: %w", err)
	}

	return nil
}

// IncrementSuccessCount increments the success count for a pattern
func (r *PostgresRepository) IncrementSuccessCount(ctx context.Context, apiSpecID uuid.UUID, patternData map[string]interface{}) error {
	// First try to find existing pattern
	query := `
		UPDATE learned_patterns
		SET success_count = success_count + 1, last_updated = $2
		WHERE api_spec_id = $1
		RETURNING id
	`

	var id uuid.UUID
	err := r.pool.QueryRow(ctx, query, apiSpecID, time.Now()).Scan(&id)
	if err != nil {
		// Pattern doesn't exist, create new one
		pattern := &entities.LearnedPattern{
			APISpecID:    apiSpecID,
			PatternData:  patternData,
			SuccessCount: 1,
		}
		return r.SaveLearnedPattern(ctx, pattern)
	}

	return nil
}

// GetLearningThreshold gets the learning threshold from system config
func (r *PostgresRepository) GetLearningThreshold(ctx context.Context) (int, error) {
	query := `SELECT value FROM system_config WHERE key = 'learning_threshold'`

	var value map[string]interface{}
	err := r.pool.QueryRow(ctx, query).Scan(&value)
	if err != nil {
		return 5, nil // Default to 5
	}

	if threshold, ok := value["value"].(float64); ok {
		return int(threshold), nil
	}

	return 5, nil
}

// CheckIfLearnedEnough checks if pattern has enough successes
func (r *PostgresRepository) CheckIfLearnedEnough(ctx context.Context, apiSpecID uuid.UUID) (bool, error) {
	threshold, err := r.GetLearningThreshold(ctx)
	if err != nil {
		threshold = 5
	}

	query := `
		SELECT success_count >= $2
		FROM learned_patterns
		WHERE api_spec_id = $1
		LIMIT 1
	`

	var learned bool
	err = r.pool.QueryRow(ctx, query, apiSpecID, threshold).Scan(&learned)
	if err != nil {
		return false, nil
	}

	return learned, nil
}

