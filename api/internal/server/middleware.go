package server

import (
	"context"
	"net/http"
	"time"

	"starterkit/internal/platform/logger"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
)

// applyMiddleware wraps the handler with all middleware
func (s *Server) applyMiddleware(h http.Handler) http.Handler {
	// Apply middleware in reverse order (innermost first)
	h = s.recoveryMiddleware(h)
	h = s.loggingMiddleware(h)
	h = s.requestIDMiddleware(h)
	h = s.corsMiddleware(h)
	return h
}

// corsMiddleware adds CORS headers to responses
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-Email, X-Request-ID")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// requestIDMiddleware adds a unique request ID to the context
func (s *Server) requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Add to response header
		w.Header().Set("X-Request-ID", requestID)

		// Add to context
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// loggingMiddleware logs HTTP requests and adds logger to context
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Get request ID from context
		requestID, _ := r.Context().Value(requestIDKey).(string)

		// Extract trace context if telemetry is enabled
		var traceID, spanID string
		if s.config.Telemetry.Enabled {
			span := trace.SpanFromContext(r.Context())
			if span.SpanContext().IsValid() {
				traceID = span.SpanContext().TraceID().String()
				spanID = span.SpanContext().SpanID().String()
			}
		}

		// Create request-specific logger
		requestLogger := s.logger.With(
			"request_id", requestID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		if traceID != "" {
			requestLogger = requestLogger.With("trace_id", traceID, "span_id", spanID)
		}

		// Add logger to context
		ctx := logger.WithContext(r.Context(), requestLogger)

		// Wrap response writer to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrapped, r.WithContext(ctx))

		// Log request completion
		requestLogger.Info("request completed",
			"status", wrapped.statusCode,
			"duration", time.Since(start),
			"bytes", wrapped.bytesWritten,
		)
	})
}

// recoveryMiddleware recovers from panics and returns 500
func (s *Server) recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				reqLogger := logger.FromContext(r.Context())
				reqLogger.Error("panic recovered",
					"error", err,
					"stack", "stack trace would go here",
				)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and bytes written
type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// RequestIDFromContext extracts the request ID from context
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}
