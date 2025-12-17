package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GeminiEmbeddingAdapter handles text embeddings via Gemini API
type GeminiEmbeddingAdapter struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// geminiEmbeddingRequest represents a Gemini embedding request
type geminiEmbeddingRequest struct {
	Model   string                 `json:"model"`
	Content geminiEmbeddingContent `json:"content"`
}

// geminiEmbeddingContent represents content for embedding
type geminiEmbeddingContent struct {
	Parts []geminiEmbeddingPart `json:"parts"`
}

// geminiEmbeddingPart represents a text part
type geminiEmbeddingPart struct {
	Text string `json:"text"`
}

// geminiEmbeddingResponse represents a Gemini embedding response
type geminiEmbeddingResponse struct {
	Embedding struct {
		Values []float32 `json:"values"`
	} `json:"embedding"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// NewGeminiEmbeddingAdapter creates a new embedding adapter using Gemini
func NewGeminiEmbeddingAdapter(apiKey string) *GeminiEmbeddingAdapter {
	return &GeminiEmbeddingAdapter{
		apiKey:     apiKey,
		model:      "text-embedding-004",
		httpClient: &http.Client{},
	}
}

// IsAvailable returns true if the adapter is configured
func (a *GeminiEmbeddingAdapter) IsAvailable() bool {
	return a.apiKey != ""
}

// GenerateEmbedding generates an embedding for a single text
func (a *GeminiEmbeddingAdapter) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if a.apiKey == "" {
		return nil, fmt.Errorf("Gemini API key not configured")
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/%s:embedContent?key=%s", a.model, a.apiKey)

	reqBody := geminiEmbeddingRequest{
		Model: fmt.Sprintf("models/%s", a.model),
		Content: geminiEmbeddingContent{
			Parts: []geminiEmbeddingPart{{Text: text}},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
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
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var embResp geminiEmbeddingResponse
	if err := json.Unmarshal(body, &embResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if embResp.Error != nil {
		return nil, fmt.Errorf("Gemini API error: %s", embResp.Error.Message)
	}

	return embResp.Embedding.Values, nil
}
