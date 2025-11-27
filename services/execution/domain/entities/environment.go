package entities

import (
	"time"

	"github.com/google/uuid"
)

// Environment represents a target environment (QA, Staging, etc.)
type Environment struct {
	ID         uuid.UUID              `json:"id"`
	Name       string                 `json:"name"`
	BaseURL    string                 `json:"base_url"`
	AuthConfig map[string]interface{} `json:"auth_config"`
	Active     bool                   `json:"active"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

// NewEnvironment creates a new environment entity
func NewEnvironment(name, baseURL string) *Environment {
	now := time.Now()
	return &Environment{
		ID:         uuid.New(),
		Name:       name,
		BaseURL:    baseURL,
		AuthConfig: make(map[string]interface{}),
		Active:     true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// Validate checks if the environment is valid
func (e *Environment) Validate() error {
	if e.Name == "" {
		return ErrInvalidEnvironment
	}
	if e.BaseURL == "" {
		return ErrInvalidEnvironment
	}
	return nil
}

