package server

import (
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// routes sets up all application routes
func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", s.handleHealthCheck())

	// API v1 routes
	v1Mux := http.NewServeMux()

	// User endpoints
	v1Mux.HandleFunc("GET /users", s.userHandler.HandleListUsers())
	v1Mux.HandleFunc("GET /users/{id}", s.userHandler.HandleGetUser())

	// Mount v1 routes
	mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1Mux))

	// Apply middleware chain
	handler := s.applyMiddleware(mux)

	// Wrap with OpenTelemetry instrumentation if enabled
	if s.config.Telemetry.Enabled {
		handler = otelhttp.NewHandler(handler, "http-server")
	}

	return handler
}
