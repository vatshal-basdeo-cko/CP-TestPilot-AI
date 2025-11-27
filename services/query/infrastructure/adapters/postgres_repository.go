package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testpilot-ai/query/domain/entities"
	"github.com/testpilot-ai/query/domain/repositories"
)

// PostgresQueryRepository implements query repository
type PostgresQueryRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresQueryRepository creates a new repository
func NewPostgresQueryRepository(pool *pgxpool.Pool) *PostgresQueryRepository {
	return &PostgresQueryRepository{pool: pool}
}

// FindExecutionByID retrieves a single execution
func (r *PostgresQueryRepository) FindExecutionByID(ctx context.Context, id uuid.UUID) (*entities.TestExecution, error) {
	query := `
		SELECT id, user_id, api_spec_id, natural_language_request,
		       constructed_request, response, validation_result,
		       status, execution_time_ms, created_at
		FROM test_executions
		WHERE id = $1
	`

	var exec entities.TestExecution
	var constructedReq, response, validationResult []byte

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&exec.ID,
		&exec.UserID,
		&exec.APISpecID,
		&exec.NaturalLanguageRequest,
		&constructedReq,
		&response,
		&validationResult,
		&exec.Status,
		&exec.ExecutionTimeMs,
		&exec.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(constructedReq, &exec.ConstructedRequest)
	json.Unmarshal(response, &exec.Response)
	json.Unmarshal(validationResult, &exec.ValidationResult)

	return &exec, nil
}

// ListExecutions retrieves executions with filters
func (r *PostgresQueryRepository) ListExecutions(ctx context.Context, filters repositories.Filters) ([]entities.TestExecution, int64, error) {
	// Build query dynamically
	where := []string{"1=1"}
	args := []interface{}{}
	argCount := 1

	if filters.UserID != nil {
		where = append(where, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, *filters.UserID)
		argCount++
	}

	if filters.APISpecID != nil {
		where = append(where, fmt.Sprintf("api_spec_id = $%d", argCount))
		args = append(args, *filters.APISpecID)
		argCount++
	}

	if filters.Status != "" {
		where = append(where, fmt.Sprintf("status = $%d", argCount))
		args = append(args, filters.Status)
		argCount++
	}

	if filters.Search != "" {
		where = append(where, fmt.Sprintf("natural_language_request ILIKE $%d", argCount))
		args = append(args, "%"+filters.Search+"%")
		argCount++
	}

	whereClause := strings.Join(where, " AND ")

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM test_executions WHERE %s", whereClause)
	var total int64
	r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)

	// Get records
	if filters.Limit == 0 {
		filters.Limit = 20
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, api_spec_id, natural_language_request,
		       constructed_request, response, validation_result,
		       status, execution_time_ms, created_at
		FROM test_executions
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	args = append(args, filters.Limit, filters.Offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var executions []entities.TestExecution
	for rows.Next() {
		var exec entities.TestExecution
		var constructedReq, response, validationResult []byte

		err := rows.Scan(
			&exec.ID,
			&exec.UserID,
			&exec.APISpecID,
			&exec.NaturalLanguageRequest,
			&constructedReq,
			&response,
			&validationResult,
			&exec.Status,
			&exec.ExecutionTimeMs,
			&exec.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal(constructedReq, &exec.ConstructedRequest)
		json.Unmarshal(response, &exec.Response)
		json.Unmarshal(validationResult, &exec.ValidationResult)

		executions = append(executions, exec)
	}

	return executions, total, nil
}

// GetAnalytics retrieves aggregated statistics
func (r *PostgresQueryRepository) GetAnalytics(ctx context.Context, startDate, endDate *time.Time) (*entities.Analytics, error) {
	analytics := &entities.Analytics{}

	// Overall stats
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as successful,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed,
			AVG(execution_time_ms) as avg_time
		FROM test_executions
		WHERE created_at >= COALESCE($1, created_at)
		  AND created_at <= COALESCE($2, created_at)
	`

	err := r.pool.QueryRow(ctx, query, startDate, endDate).Scan(
		&analytics.TotalTests,
		&analytics.SuccessfulTests,
		&analytics.FailedTests,
		&analytics.AvgExecutionTime,
	)
	if err != nil {
		return nil, err
	}

	if analytics.TotalTests > 0 {
		analytics.SuccessRate = float64(analytics.SuccessfulTests) / float64(analytics.TotalTests) * 100
	}

	return analytics, nil
}

// GetAPIAnalytics retrieves per-API statistics
func (r *PostgresQueryRepository) GetAPIAnalytics(ctx context.Context, apiSpecID uuid.UUID) (*entities.APIStats, error) {
	query := `
		SELECT 
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'success' THEN 1 END) as successful
		FROM test_executions
		WHERE api_spec_id = $1
	`

	stats := &entities.APIStats{}
	err := r.pool.QueryRow(ctx, query, apiSpecID).Scan(
		&stats.TestCount,
		&stats.SuccessCount,
	)
	if err != nil {
		return nil, err
	}

	if stats.TestCount > 0 {
		stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.TestCount) * 100
	}

	return stats, nil
}

