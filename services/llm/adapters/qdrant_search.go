package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/testpilot-ai/llm/domain/entities"
)

// QdrantSearchAdapter handles vector search for RAG
type QdrantSearchAdapter struct {
	baseURL    string
	collection string
	httpClient *http.Client
}

// NewQdrantSearchAdapter creates a new Qdrant search adapter
func NewQdrantSearchAdapter(baseURL, collection string) *QdrantSearchAdapter {
	return &QdrantSearchAdapter{
		baseURL:    baseURL,
		collection: collection,
		httpClient: &http.Client{},
	}
}

// SearchRequest represents a Qdrant search request
type SearchRequest struct {
	Vector      []float32 `json:"vector"`
	Limit       int       `json:"limit"`
	WithPayload bool      `json:"with_payload"`
}

// SearchResult represents a Qdrant search result
type SearchResult struct {
	ID      string                 `json:"id"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

// Search performs a vector similarity search
func (a *QdrantSearchAdapter) Search(vector []float32, limit int) ([]entities.RetrievalContext, error) {
	reqBody := SearchRequest{
		Vector:      vector,
		Limit:       limit,
		WithPayload: true,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%s/points/search", a.baseURL, a.collection)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
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
		Result []SearchResult `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to RetrievalContext
	var contexts []entities.RetrievalContext
	for _, r := range result.Result {
		ctx := entities.RetrievalContext{
			Score: r.Score,
		}

		// Extract fields from payload
		if name, ok := r.Payload["api_name"].(string); ok {
			ctx.APIName = name
		}
		if version, ok := r.Payload["version"].(string); ok {
			ctx.Version = version
		}
		if desc, ok := r.Payload["description"].(string); ok {
			ctx.Description = desc
		}
		if config, ok := r.Payload["config"].(string); ok {
			_ = json.Unmarshal([]byte(config), &ctx.Config)
		}

		contexts = append(contexts, ctx)
	}

	return contexts, nil
}

// SearchByText converts text to embedding and searches
func (a *QdrantSearchAdapter) SearchByText(text string, embedding []float32, limit int) ([]entities.RetrievalContext, error) {
	return a.Search(embedding, limit)
}

