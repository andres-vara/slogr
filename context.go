package slogr

import (
	"context"
)

// ContextKey is a type used for context keys to avoid collisions
type ContextKey string

const (
	// LoggerKey is the context key for storing the logger
	LoggerKey ContextKey = "slogr_logger"
)

// WithLogger returns a new context with the logger added
func WithLogger(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// FromContext extracts a logger from the context
// Returns nil if no logger is found
func FromContext(ctx context.Context) *Logger {
	if logger, ok := ctx.Value(LoggerKey).(*Logger); ok {
		return logger
	}
	return nil
}