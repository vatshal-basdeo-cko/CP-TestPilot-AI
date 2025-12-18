package logger

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var (
	// globalLogger is the global logger instance
	globalLogger zerolog.Logger
	// serviceName is the name of the service using this logger
	serviceName string
)

// Init initializes the global logger with the given service name
func Init(name string) {
	serviceName = name
	
	// Set log level from environment variable
	logLevel := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	if logLevel == "" {
		logLevel = "INFO"
	}

	// Parse log level
	var level zerolog.Level
	switch logLevel {
	case "DEBUG":
		level = zerolog.DebugLevel
	case "INFO":
		level = zerolog.InfoLevel
	case "ERROR":
		level = zerolog.ErrorLevel
	default:
		level = zerolog.InfoLevel
	}

	// Configure zerolog for human-readable text output
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = time.RFC3339

	// Create console writer with human-readable format
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    true, // No color for Docker logs
	}

	// Build logger with service name
	globalLogger = zerolog.New(output).
		With().
		Str("service", serviceName).
		Timestamp().
		Logger()
}

// Logger returns the global logger instance
func Logger() *zerolog.Logger {
	return &globalLogger
}

// WithContext returns a logger with context fields (request ID, etc.)
func WithContext(ctx context.Context) zerolog.Logger {
	log := globalLogger.With()
	
	// Extract request ID from context if available
	if requestID, ok := ctx.Value("request_id").(string); ok && requestID != "" {
		log = log.Str("request_id", requestID)
	}
	
	return log.Logger()
}

// Debug logs a debug message
func Debug(msg string) {
	globalLogger.Debug().Msg(msg)
}

// Debugf logs a formatted debug message
func Debugf(format string, v ...interface{}) {
	globalLogger.Debug().Msgf(format, v...)
}

// Info logs an info message
func Info(msg string) {
	globalLogger.Info().Msg(msg)
}

// Infof logs a formatted info message
func Infof(format string, v ...interface{}) {
	globalLogger.Info().Msgf(format, v...)
}

// Error logs an error message
func Error(msg string) {
	globalLogger.Error().Msg(msg)
}

// Errorf logs a formatted error message
func Errorf(format string, v ...interface{}) {
	globalLogger.Error().Msgf(format, v...)
}

// Warn logs a warning message
func Warn(msg string) {
	globalLogger.Warn().Msg(msg)
}

// Warnf logs a formatted warning message
func Warnf(format string, v ...interface{}) {
	globalLogger.Warn().Msgf(format, v...)
}

// Err logs an error with an error object
func Err(err error) *zerolog.Event {
	return globalLogger.Error().Err(err)
}

// WithRequestID creates a logger with a request ID
func WithRequestID(requestID string) zerolog.Logger {
	return globalLogger.With().Str("request_id", requestID).Logger()
}

// WithFields creates a logger with additional fields
func WithFields(fields map[string]interface{}) zerolog.Logger {
	log := globalLogger.With()
	for k, v := range fields {
		log = log.Interface(k, v)
	}
	return log.Logger()
}

