package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// QdrantAdapter handles vector database operations
type QdrantAdapter struct {
	baseURL    string
	collection string
	httpClient *http.Client
}

// QdrantPoint represents a point in Qdrant
type QdrantPoint struct {
	ID      string                 `json:"id"`
	Vector  []float32              `json:"vector"`
	Payload map[string]interface{} `json:"payload"`
}

// QdrantUpsertRequest represents an upsert request
type QdrantUpsertRequest struct {
	Points []QdrantPoint `json:"points"`
}

// QdrantSearchRequest represents a search request
type QdrantSearchRequest struct {
	Vector      []float32              `json:"vector"`
	Limit       int                    `json:"limit"`
	WithPayload bool                   `json:"with_payload"`
	Filter      map[string]interface{} `json:"filter,omitempty"`
}

// QdrantSearchResult represents a search result
type QdrantSearchResult struct {
	ID      string                 `json:"id"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

// NewQdrantAdapter creates a new Qdrant adapter
func NewQdrantAdapter(baseURL, collection string) *QdrantAdapter {
	return &QdrantAdapter{
		baseURL:    baseURL,
		collection: collection,
		httpClient: &http.Client{},
	}
}

// EnsureCollection ensures the collection exists
func (a *QdrantAdapter) EnsureCollection(vectorSize int) error {
	// Check if collection exists
	resp, err := a.httpClient.Get(a.baseURL + "/collections/" + a.collection)
	if err != nil {
		return fmt.Errorf("failed to check collection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil // Collection exists
	}

	// Create collection
	createReq := map[string]interface{}{
		"vectors": map[string]interface{}{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}

	jsonBody, err := json.Marshal(createReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("PUT", a.baseURL+"/collections/"+a.collection, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create collection (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// Upsert inserts or updates points in the collection
func (a *QdrantAdapter) Upsert(points []QdrantPoint) error {
	reqBody := QdrantUpsertRequest{Points: points}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("PUT", a.baseURL+"/collections/"+a.collection+"/points", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upsert points: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to upsert points (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// Search searches for similar vectors
func (a *QdrantAdapter) Search(vector []float32, limit int) ([]QdrantSearchResult, error) {
	reqBody := QdrantSearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", a.baseURL+"/collections/"+a.collection+"/points/search", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Result []QdrantSearchResult `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Result, nil
}

// Delete deletes a point by ID
func (a *QdrantAdapter) Delete(id uuid.UUID) error {
	reqBody := map[string]interface{}{
		"points": []string{id.String()},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", a.baseURL+"/collections/"+a.collection+"/points/delete", bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete point: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete point (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

