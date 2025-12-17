package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// EmbeddingService handles text embeddings via Gemini API
type EmbeddingService struct {
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

// NewEmbeddingService creates a new embedding service using Gemini
func NewEmbeddingService(apiKey string) *EmbeddingService {
	return &EmbeddingService{
		apiKey:     apiKey,
		model:      "text-embedding-004",
		httpClient: &http.Client{},
	}
}

// GenerateEmbedding generates an embedding for a single text
func (s *EmbeddingService) GenerateEmbedding(text string) ([]float32, error) {
	if s.apiKey == "" {
		// Return zero vectors if no API key (for testing)
		return make([]float32, 768), nil
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/%s:embedContent?key=%s", s.model, s.apiKey)

	reqBody := geminiEmbeddingRequest{
		Model: fmt.Sprintf("models/%s", s.model),
		Content: geminiEmbeddingContent{
			Parts: []geminiEmbeddingPart{{Text: text}},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
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

// GenerateEmbeddings generates embeddings for multiple texts
func (s *EmbeddingService) GenerateEmbeddings(texts []string) ([][]float32, error) {
	result := make([][]float32, len(texts))

	for i, text := range texts {
		embedding, err := s.GenerateEmbedding(text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}
		result[i] = embedding
	}

	return result, nil
}
