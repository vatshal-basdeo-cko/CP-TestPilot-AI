package adapters

import (
	"context"
)

// LLMProvider interface for different LLM providers
type LLMProvider interface {
	// Complete sends a prompt and returns the completion
	Complete(ctx context.Context, prompt string) (string, error)
	
	// Name returns the provider name
	Name() string
	
	// IsAvailable checks if the provider is configured
	IsAvailable() bool
}

// LLMMessage represents a chat message
type LLMMessage struct {
	Role    string `json:"role"`    // system, user, assistant
	Content string `json:"content"`
}

