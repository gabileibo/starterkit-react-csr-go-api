package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Service   ServiceConfig
	Server    ServerConfig
	Database  DatabaseConfig
	Telemetry TelemetryConfig
}

// ServiceConfig contains service metadata
type ServiceConfig struct {
	Name    string
	Version string
}

// ServerConfig contains HTTP server configuration
type ServerConfig struct {
	Address         string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// DatabaseConfig contains database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// TelemetryConfig contains observability configuration
type TelemetryConfig struct {
	OTLPEndpoint string
	Enabled      bool
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// It's ok if .env doesn't exist in production
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	cfg := &Config{
		Service: ServiceConfig{
			Name:    getEnv("SERVICE_NAME", "starterkit"),
			Version: getEnv("SERVICE_VERSION", "1.0.0"),
		},
		Server: ServerConfig{
			Address:         getEnv("SERVER_ADDRESS", ":8080"),
			ReadTimeout:     getDuration("SERVER_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getDuration("SERVER_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second),
			ShutdownTimeout: getDuration("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			Database:        getEnv("DB_NAME", "starterkit"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDuration("DB_CONN_MAX_IDLE_TIME", 1*time.Minute),
		},
		Telemetry: TelemetryConfig{
			OTLPEndpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			Enabled:      getBoolEnv("TELEMETRY_ENABLED", true),
		},
	}

	return cfg, nil
}

// DSN returns the PostgreSQL connection string
func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intValue int
		if _, err := fmt.Sscanf(value, "%d", &intValue); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

func getDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
