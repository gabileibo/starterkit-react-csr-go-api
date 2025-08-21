package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"starterkit/internal/config"
	"starterkit/internal/db"
	"starterkit/internal/users"
)

// Server represents the HTTP server
type Server struct {
	httpServer  *http.Server
	config      *config.Config
	logger      *slog.Logger
	queries     *db.Queries
	userHandler *users.Handler
}

// New creates a new server instance
func New(cfg *config.Config, logger *slog.Logger, queries *db.Queries) *Server {
	// Create services
	userService := users.NewService(queries)

	// Create handlers
	userHandler := users.NewHandler(userService, logger)

	s := &Server{
		config:      cfg,
		logger:      logger,
		queries:     queries,
		userHandler: userHandler,
	}

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      s.routes(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	return s
}

// Start begins listening for HTTP requests
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// handleHealthCheck returns a simple health check handler
func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"healthy","service":"%s","version":"%s"}`,
			s.config.Service.Name, s.config.Service.Version)
	}
}
