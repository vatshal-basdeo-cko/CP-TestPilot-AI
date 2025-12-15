package adapters

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements LLMProvider for OpenAI
type OpenAIProvider struct {
	client *openai.Client
	model  string
	apiKey string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) *OpenAIProvider {
	var client *openai.Client
	if apiKey != "" {
		client = openai.NewClient(apiKey)
	}
	return &OpenAIProvider{
		client: client,
		model:  openai.GPT4TurboPreview,
		apiKey: apiKey,
	}
}

// Complete sends a prompt to OpenAI and returns the completion
func (p *OpenAIProvider) Complete(ctx context.Context, prompt string) (string, error) {
	if p.client == nil {
		return "", fmt.Errorf("OpenAI client not initialized")
	}

	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
			MaxTokens:   2000,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// Name returns the provider name
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// IsAvailable checks if OpenAI is configured
func (p *OpenAIProvider) IsAvailable() bool {
	return p.apiKey != "" && p.client != nil
}

// GenerateEmbedding generates embeddings using OpenAI
func (p *OpenAIProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if p.client == nil {
		return nil, fmt.Errorf("OpenAI client not initialized")
	}

	resp, err := p.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequest{
			Model: openai.AdaEmbeddingV2,
			Input: []string{text},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("OpenAI embedding error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return resp.Data[0].Embedding, nil
}

