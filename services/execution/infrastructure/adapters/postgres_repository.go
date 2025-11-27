package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testpilot-ai/execution/domain/entities"
)

// PostgresRepository implements execution repository using PostgreSQL
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

// SaveExecution saves an execution record
func (r *PostgresRepository) SaveExecution(ctx context.Context, request *entities.APIRequest, response *entities.APIResponse) error {
	query := `
		INSERT INTO test_executions (
			id, user_id, api_spec_id, natural_language_request,
			constructed_request, response, validation_result,
			status, execution_time_ms, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Marshal request and response to JSON
	constructedReq, err := json.Marshal(map[string]interface{}{
		"method":        request.Method,
		"url":           request.URL,
		"headers":       request.Headers,
		"query_params":  request.QueryParams,
		"body":          request.Body,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	responseJSON, err := json.Marshal(map[string]interface{}{
		"status_code":       response.StatusCode,
		"headers":           response.Headers,
		"body":              response.Body,
		"execution_time_ms": response.ExecutionTimeMs,
		"error":             response.Error,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	status := "success"
	if !response.Success {
		status = "failed"
	}

	_, err = r.pool.Exec(ctx, query,
		response.ID,
		nil, // user_id - to be set when auth is implemented
		request.APISpecID,
		"", // natural_language_request - from LLM service
		constructedReq,
		responseJSON,
		nil, // validation_result - from validation service
		status,
		response.ExecutionTimeMs,
		time.Now(),
	)

	return err
}

// FindExecutionByID retrieves an execution by ID
func (r *PostgresRepository) FindExecutionByID(ctx context.Context, id uuid.UUID) (*entities.APIResponse, error) {
	query := `
		SELECT id, constructed_request, response, status, execution_time_ms, created_at
		FROM test_executions
		WHERE id = $1
	`

	var response entities.APIResponse
	var constructedReq, responseJSON []byte
	var status string
	var createdAt time.Time

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&response.ID,
		&constructedReq,
		&responseJSON,
		&status,
		&response.ExecutionTimeMs,
		&createdAt,
	)
	if err != nil {
		return nil, err
	}

	// Parse response JSON
	var respData map[string]interface{}
	if err := json.Unmarshal(responseJSON, &respData); err == nil {
		response.StatusCode = int(respData["status_code"].(float64))
		response.Body = respData["body"]
		if errMsg, ok := respData["error"].(string); ok {
			response.Error = errMsg
		}
	}

	response.Success = status == "success"
	response.Timestamp = createdAt

	return &response, nil
}

// ListExecutions retrieves executions with pagination
func (r *PostgresRepository) ListExecutions(ctx context.Context, limit, offset int) ([]*entities.APIResponse, error) {
	query := `
		SELECT id, constructed_request, response, status, execution_time_ms, created_at
		FROM test_executions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var executions []*entities.APIResponse
	for rows.Next() {
		var response entities.APIResponse
		var constructedReq, responseJSON []byte
		var status string
		var createdAt time.Time

		err := rows.Scan(
			&response.ID,
			&constructedReq,
			&responseJSON,
			&status,
			&response.ExecutionTimeMs,
			&createdAt,
		)
		if err != nil {
			continue
		}

		// Parse response JSON
		var respData map[string]interface{}
		if err := json.Unmarshal(responseJSON, &respData); err == nil {
			if statusCode, ok := respData["status_code"].(float64); ok {
				response.StatusCode = int(statusCode)
			}
			response.Body = respData["body"]
		}

		response.Success = status == "success"
		response.Timestamp = createdAt

		executions = append(executions, &response)
	}

	return executions, nil
}

