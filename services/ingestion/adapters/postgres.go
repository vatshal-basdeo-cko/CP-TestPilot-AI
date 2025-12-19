package adapters

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testpilot-ai/ingestion/domain/entities"
)

// PostgresRepository handles database operations
type PostgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{pool: pool}
}

// SaveAPISpecification saves or updates an API specification
func (r *PostgresRepository) SaveAPISpecification(ctx context.Context, spec *entities.APISpecification) error {
	query := `
		INSERT INTO api_specifications (id, name, version, source_type, source_path, content_hash, metadata, created_at, updated_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			version = EXCLUDED.version,
			content_hash = EXCLUDED.content_hash,
			metadata = EXCLUDED.metadata,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(ctx, query,
		spec.ID,
		spec.Name,
		spec.Version,
		spec.SourceType,
		spec.SourcePath,
		spec.ContentHash,
		spec.Metadata,
		spec.CreatedAt,
		spec.UpdatedAt,
		spec.CreatedBy,
	)

	if err != nil {
		return fmt.Errorf("failed to save API specification: %w", err)
	}

	return nil
}

// GetAPISpecificationByHash retrieves an API specification by content hash
func (r *PostgresRepository) GetAPISpecificationByHash(ctx context.Context, contentHash string) (*entities.APISpecification, error) {
	query := `
		SELECT id, name, version, source_type, source_path, content_hash, metadata, created_at, updated_at, created_by
		FROM api_specifications
		WHERE content_hash = $1
		LIMIT 1
	`

	var spec entities.APISpecification
	err := r.pool.QueryRow(ctx, query, contentHash).Scan(
		&spec.ID,
		&spec.Name,
		&spec.Version,
		&spec.SourceType,
		&spec.SourcePath,
		&spec.ContentHash,
		&spec.Metadata,
		&spec.CreatedAt,
		&spec.UpdatedAt,
		&spec.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &spec, nil
}

// GetAPISpecificationByNameVersion retrieves an API specification by name and version
func (r *PostgresRepository) GetAPISpecificationByNameVersion(ctx context.Context, name, version string) (*entities.APISpecification, error) {
	query := `
		SELECT id, name, version, source_type, source_path, content_hash, metadata, created_at, updated_at, created_by
		FROM api_specifications
		WHERE name = $1 AND version = $2
		LIMIT 1
	`

	var spec entities.APISpecification
	err := r.pool.QueryRow(ctx, query, name, version).Scan(
		&spec.ID,
		&spec.Name,
		&spec.Version,
		&spec.SourceType,
		&spec.SourcePath,
		&spec.ContentHash,
		&spec.Metadata,
		&spec.CreatedAt,
		&spec.UpdatedAt,
		&spec.CreatedBy,
	)

	if err != nil {
		return nil, err
	}

	return &spec, nil
}

// UpdateAPISpecification updates an existing API specification
func (r *PostgresRepository) UpdateAPISpecification(ctx context.Context, spec *entities.APISpecification) error {
	query := `
		UPDATE api_specifications 
		SET content_hash = $1, metadata = $2, updated_at = $3, source_path = $4
		WHERE id = $5
	`

	_, err := r.pool.Exec(ctx, query,
		spec.ContentHash,
		spec.Metadata,
		spec.UpdatedAt,
		spec.SourcePath,
		spec.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update API specification: %w", err)
	}

	return nil
}

// GetAllAPISpecifications retrieves all API specifications
func (r *PostgresRepository) GetAllAPISpecifications(ctx context.Context) ([]entities.APISpecification, error) {
	query := `
		SELECT id, name, version, source_type, source_path, content_hash, metadata, created_at, updated_at, created_by
		FROM api_specifications
		ORDER BY updated_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query API specifications: %w", err)
	}
	defer rows.Close()

	var specs []entities.APISpecification
	for rows.Next() {
		var spec entities.APISpecification
		err := rows.Scan(
			&spec.ID,
			&spec.Name,
			&spec.Version,
			&spec.SourceType,
			&spec.SourcePath,
			&spec.ContentHash,
			&spec.Metadata,
			&spec.CreatedAt,
			&spec.UpdatedAt,
			&spec.CreatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		specs = append(specs, spec)
	}

	return specs, nil
}

// DeleteAPISpecification deletes an API specification
func (r *PostgresRepository) DeleteAPISpecification(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM api_specifications WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete API specification: %w", err)
	}
	return nil
}

// SaveIngestionLog saves an ingestion log entry
func (r *PostgresRepository) SaveIngestionLog(ctx context.Context, result *entities.IngestionResult) error {
	query := `
		INSERT INTO ingestion_logs (id, source_type, source_path, status, apis_ingested, error_message, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.pool.Exec(ctx, query,
		result.ID,
		result.SourceType,
		result.SourcePath,
		result.Status,
		result.APIsIngested,
		result.ErrorMessage,
		result.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to save ingestion log: %w", err)
	}

	return nil
}

// GetIngestionLogs retrieves recent ingestion logs
func (r *PostgresRepository) GetIngestionLogs(ctx context.Context, limit int) ([]entities.IngestionResult, error) {
	query := `
		SELECT id, source_type, source_path, status, apis_ingested, error_message, created_at
		FROM ingestion_logs
		ORDER BY created_at DESC
		LIMIT $1
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query ingestion logs: %w", err)
	}
	defer rows.Close()

	var logs []entities.IngestionResult
	for rows.Next() {
		var log entities.IngestionResult
		var errorMsg *string
		err := rows.Scan(
			&log.ID,
			&log.SourceType,
			&log.SourcePath,
			&log.Status,
			&log.APIsIngested,
			&errorMsg,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if errorMsg != nil {
			log.ErrorMessage = *errorMsg
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// NewIngestionResult creates a new ingestion result
func NewIngestionResult(sourceType, sourcePath, status string, apisIngested int, errorMessage string) *entities.IngestionResult {
	return &entities.IngestionResult{
		ID:           uuid.New(),
		SourceType:   sourceType,
		SourcePath:   sourcePath,
		Status:       status,
		APIsIngested: apisIngested,
		ErrorMessage: errorMessage,
		CreatedAt:    time.Now(),
	}
}

