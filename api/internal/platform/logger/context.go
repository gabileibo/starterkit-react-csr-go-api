package logger

import (
	"context"
	"log/slog"
)

type contextKey string

const (
	loggerKey contextKey = "logger"
)

// WithContext adds a logger to the context
func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// FromContext extracts the logger from context
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}
