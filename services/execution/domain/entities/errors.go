package entities

import "errors"

var (
	ErrInvalidMethod      = errors.New("invalid HTTP method")
	ErrInvalidURL         = errors.New("invalid URL")
	ErrInvalidEnvironment = errors.New("invalid environment")
	ErrExecutionFailed    = errors.New("execution failed")
	ErrTimeout            = errors.New("request timeout")
	ErrEnvironmentNotFound = errors.New("environment not found")
)

