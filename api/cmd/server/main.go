package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"starterkit/internal/config"
	"starterkit/internal/db"
	"starterkit/internal/platform/database"
	"starterkit/internal/platform/telemetry"
	"starterkit/internal/server"
)

func main() {
	// Initialize structured logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize telemetry
	shutdown, err := telemetry.Init(context.Background(), cfg.Service.Name, cfg.Service.Version)
	if err != nil {
		logger.Error("failed to initialize telemetry", "error", err)
		os.Exit(1)
	}
	defer shutdown()

	// Initialize database connection
	dbPool, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	// Initialize sqlc queries
	queries := db.New(dbPool)

	// Initialize server
	srv := server.New(cfg, logger, queries)

	// Start server in a goroutine
	go func() {
		logger.Info("starting server", "address", cfg.Server.Address)
		if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", "error", err)
		os.Exit(1)
	}

	logger.Info("server exited")
}
