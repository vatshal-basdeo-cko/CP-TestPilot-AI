package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/testpilot-ai/ingestion/adapters"
	"github.com/testpilot-ai/ingestion/domain/entities"
)

// IngestionHandler handles ingestion-related HTTP requests
type IngestionHandler struct {
	fileParser     *adapters.FileParser
	postmanParser  *adapters.PostmanParser
	embeddingService *adapters.EmbeddingService
	qdrantAdapter  *adapters.QdrantAdapter
	postgresRepo   *adapters.PostgresRepository
}

// NewIngestionHandler creates a new ingestion handler
func NewIngestionHandler(
	fileParser *adapters.FileParser,
	postmanParser *adapters.PostmanParser,
	embeddingService *adapters.EmbeddingService,
	qdrantAdapter *adapters.QdrantAdapter,
	postgresRepo *adapters.PostgresRepository,
) *IngestionHandler {
	return &IngestionHandler{
		fileParser:     fileParser,
		postmanParser:  postmanParser,
		embeddingService: embeddingService,
		qdrantAdapter:  qdrantAdapter,
		postgresRepo:   postgresRepo,
	}
}

// Health returns service health status
func (h *IngestionHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "ingestion",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// IngestFile handles single file ingestion
func (h *IngestionHandler) IngestFile(c *gin.Context) {
	var req struct {
		FilePath string `json:"file_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_path is required"})
		return
	}

	// Parse the file
	config, contentHash, err := h.fileParser.ParseFile(req.FilePath)
	if err != nil {
		h.logIngestion(c, "file", req.FilePath, "failed", 0, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse file: %s", err)})
		return
	}

	// Check if already ingested (same hash)
	existing, _ := h.postgresRepo.GetAPISpecificationByHash(c.Request.Context(), contentHash)
	if existing != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "File already ingested (no changes detected)",
			"api_id":  existing.ID,
		})
		return
	}

	// Check if same name+version exists (update scenario)
	existingByName, _ := h.postgresRepo.GetAPISpecificationByNameVersion(c.Request.Context(), config.Name, config.Version)
	if existingByName != nil {
		// Update existing: delete old Qdrant vector, then update
		apiID, err := h.updateExistingSpec(c, existingByName, config, contentHash, "file", req.FilePath)
		if err != nil {
			h.logIngestion(c, "file", req.FilePath, "failed", 0, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %s", err)})
			return
		}

		h.logIngestion(c, "file", req.FilePath, "updated", 1, "")
		c.JSON(http.StatusOK, gin.H{
			"message": "File updated successfully",
			"api_id":  apiID,
			"name":    config.Name,
			"version": config.Version,
		})
		return
	}

	// Generate embeddings and store (new file)
	apiID, err := h.processAndStore(c, config, contentHash, "file", req.FilePath)
	if err != nil {
		h.logIngestion(c, "file", req.FilePath, "failed", 0, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process: %s", err)})
		return
	}

	h.logIngestion(c, "file", req.FilePath, "success", 1, "")
	c.JSON(http.StatusOK, gin.H{
		"message": "File ingested successfully",
		"api_id":  apiID,
		"name":    config.Name,
		"version": config.Version,
	})
}

// IngestFolder handles folder scanning and ingestion
func (h *IngestionHandler) IngestFolder(c *gin.Context) {
	var req struct {
		FolderPath string `json:"folder_path" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folder_path is required"})
		return
	}

	// Scan folder for files
	files, err := h.fileParser.ScanFolder(req.FolderPath)
	if err != nil {
		h.logIngestion(c, "folder", req.FolderPath, "failed", 0, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to scan folder: %s", err)})
		return
	}

	// Process each file
	var ingested, skipped, failed int
	var errors []string

	for _, filePath := range files {
		config, contentHash, err := h.fileParser.ParseFile(filePath)
		if err != nil {
			failed++
			errors = append(errors, fmt.Sprintf("%s: %s", filePath, err))
			continue
		}

		// Check if already ingested
		existing, _ := h.postgresRepo.GetAPISpecificationByHash(c.Request.Context(), contentHash)
		if existing != nil {
			skipped++
			continue
		}

		// Process and store
		_, err = h.processAndStore(c, config, contentHash, "file", filePath)
		if err != nil {
			failed++
			errors = append(errors, fmt.Sprintf("%s: %s", filePath, err))
			continue
		}

		ingested++
	}

	status := "success"
	if failed > 0 {
		status = "partial"
	}
	h.logIngestion(c, "folder", req.FolderPath, status, ingested, "")

	c.JSON(http.StatusOK, gin.H{
		"message":  "Folder scan complete",
		"ingested": ingested,
		"skipped":  skipped,
		"failed":   failed,
		"errors":   errors,
	})
}

// IngestPostman handles Postman collection upload
func (h *IngestionHandler) IngestPostman(c *gin.Context) {
	// Get file from multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Read file content
	content := make([]byte, header.Size)
	_, err = file.Read(content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}

	// Parse Postman collection
	config, contentHash, err := h.postmanParser.ParseCollectionData(content)
	if err != nil {
		h.logIngestion(c, "postman", header.Filename, "failed", 0, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to parse collection: %s", err)})
		return
	}

	// Check if already ingested (same hash = no changes)
	existing, _ := h.postgresRepo.GetAPISpecificationByHash(c.Request.Context(), contentHash)
	if existing != nil {
		c.JSON(http.StatusOK, gin.H{
			"message":   "Collection already ingested (no changes detected)",
			"api_id":    existing.ID,
			"name":      config.Name,
			"endpoints": len(config.Endpoints),
		})
		return
	}

	// Check if same name+version exists (update scenario)
	existingByName, _ := h.postgresRepo.GetAPISpecificationByNameVersion(c.Request.Context(), config.Name, config.Version)
	if existingByName != nil {
		// Update existing: delete old Qdrant vector, then update
		apiID, err := h.updateExistingSpec(c, existingByName, config, contentHash, "postman", header.Filename)
		if err != nil {
			h.logIngestion(c, "postman", header.Filename, "failed", 0, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update: %s", err)})
			return
		}

		h.logIngestion(c, "postman", header.Filename, "updated", 1, "")
		c.JSON(http.StatusOK, gin.H{
			"message":   "Postman collection updated successfully",
			"api_id":    apiID,
			"name":      config.Name,
			"endpoints": len(config.Endpoints),
		})
		return
	}

	// Process and store (new collection)
	apiID, err := h.processAndStore(c, config, contentHash, "postman", header.Filename)
	if err != nil {
		h.logIngestion(c, "postman", header.Filename, "failed", 0, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to process: %s", err)})
		return
	}

	h.logIngestion(c, "postman", header.Filename, "success", 1, "")
	c.JSON(http.StatusOK, gin.H{
		"message":   "Postman collection ingested successfully",
		"api_id":    apiID,
		"name":      config.Name,
		"endpoints": len(config.Endpoints),
	})
}

// GetStatus returns ingestion status and logs
func (h *IngestionHandler) GetStatus(c *gin.Context) {
	logs, err := h.postgresRepo.GetIngestionLogs(c.Request.Context(), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "ready",
		"recent_logs": logs,
	})
}

// ListAPIs returns all ingested APIs
func (h *IngestionHandler) ListAPIs(c *gin.Context) {
	specs, err := h.postgresRepo.GetAllAPISpecifications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list APIs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"apis":  specs,
		"count": len(specs),
	})
}

// DeleteAPI deletes an API specification and its Qdrant vectors
func (h *IngestionHandler) DeleteAPI(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid API ID format"})
		return
	}

	// Delete from Qdrant first
	if err := h.qdrantAdapter.Delete(id); err != nil {
		// Log but continue - Qdrant vector may not exist
		fmt.Printf("Warning: failed to delete from Qdrant for %s: %v\n", idStr, err)
	}

	// Delete from PostgreSQL
	if err := h.postgresRepo.DeleteAPISpecification(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete API: %s", err)})
		return
	}

	h.logIngestion(c, "delete", idStr, "success", 1, "")
	c.JSON(http.StatusOK, gin.H{
		"message": "API deleted successfully",
		"id":      idStr,
	})
}

// processAndStore generates embeddings and stores the API config
func (h *IngestionHandler) processAndStore(c *gin.Context, config *entities.APIConfig, contentHash, sourceType, sourcePath string) (uuid.UUID, error) {
	apiID := uuid.New()
	now := time.Now()

	// Generate text for embedding
	embeddingText := h.generateEmbeddingText(config)

	// Generate embedding
	embedding, err := h.embeddingService.GenerateEmbedding(embeddingText)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Store in Qdrant
	configJSON, _ := json.Marshal(config)
	point := adapters.QdrantPoint{
		ID:     apiID.String(),
		Vector: embedding,
		Payload: map[string]interface{}{
			"api_name":    config.Name,
			"version":     config.Version,
			"description": config.Description,
			"endpoints":   len(config.Endpoints),
			"config":      string(configJSON),
		},
	}

	if err := h.qdrantAdapter.Upsert([]adapters.QdrantPoint{point}); err != nil {
		return uuid.Nil, fmt.Errorf("failed to store in Qdrant: %w", err)
	}

	// Store metadata in PostgreSQL
	spec := &entities.APISpecification{
		ID:          apiID,
		Name:        config.Name,
		Version:     config.Version,
		SourceType:  sourceType,
		SourcePath:  sourcePath,
		ContentHash: contentHash,
		Metadata: map[string]interface{}{
			"description": config.Description,
			"base_url":    config.BaseURL,
			"endpoints":   len(config.Endpoints),
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.postgresRepo.SaveAPISpecification(c.Request.Context(), spec); err != nil {
		return uuid.Nil, fmt.Errorf("failed to save to database: %w", err)
	}

	return apiID, nil
}

// generateEmbeddingText generates text for embedding from API config
func (h *IngestionHandler) generateEmbeddingText(config *entities.APIConfig) string {
	text := fmt.Sprintf("API: %s\nVersion: %s\nDescription: %s\n\nEndpoints:\n",
		config.Name, config.Version, config.Description)

	for _, ep := range config.Endpoints {
		text += fmt.Sprintf("- %s %s: %s\n", ep.Method, ep.Path, ep.Description)
		for _, p := range ep.Parameters {
			text += fmt.Sprintf("  Parameter: %s (%s) - %s\n", p.Name, p.Type, p.Description)
		}
	}

	return text
}

// updateExistingSpec updates an existing API specification with new content
func (h *IngestionHandler) updateExistingSpec(c *gin.Context, existing *entities.APISpecification, config *entities.APIConfig, contentHash, sourceType, sourcePath string) (uuid.UUID, error) {
	now := time.Now()

	// Delete old Qdrant vector
	if err := h.qdrantAdapter.Delete(existing.ID); err != nil {
		// Log but continue - old vector may not exist
		fmt.Printf("Warning: failed to delete old Qdrant vector for %s: %v\n", existing.ID, err)
	}

	// Generate new embedding text
	embeddingText := h.generateEmbeddingText(config)

	// Generate new embedding
	embedding, err := h.embeddingService.GenerateEmbedding(embeddingText)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Store new vector in Qdrant with same ID
	configJSON, _ := json.Marshal(config)
	point := adapters.QdrantPoint{
		ID:     existing.ID.String(),
		Vector: embedding,
		Payload: map[string]interface{}{
			"api_name":    config.Name,
			"version":     config.Version,
			"description": config.Description,
			"endpoints":   len(config.Endpoints),
			"config":      string(configJSON),
		},
	}

	if err := h.qdrantAdapter.Upsert([]adapters.QdrantPoint{point}); err != nil {
		return uuid.Nil, fmt.Errorf("failed to store in Qdrant: %w", err)
	}

	// Update PostgreSQL record
	existing.ContentHash = contentHash
	existing.SourcePath = sourcePath
	existing.UpdatedAt = now
	existing.Metadata = map[string]interface{}{
		"description": config.Description,
		"base_url":    config.BaseURL,
		"endpoints":   len(config.Endpoints),
	}

	if err := h.postgresRepo.UpdateAPISpecification(c.Request.Context(), existing); err != nil {
		return uuid.Nil, fmt.Errorf("failed to update database: %w", err)
	}

	return existing.ID, nil
}

// logIngestion logs an ingestion operation
func (h *IngestionHandler) logIngestion(c *gin.Context, sourceType, sourcePath, status string, apisIngested int, errorMessage string) {
	result := adapters.NewIngestionResult(sourceType, sourcePath, status, apisIngested, errorMessage)
	_ = h.postgresRepo.SaveIngestionLog(c.Request.Context(), result)
}

