// Package config loads runtime configuration from environment variables.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration for the API server.
type Config struct {
	HTTPPort        string
	CORSAllowOrigin string

	DBDriver string // "sqlite" for MVP; "postgres"/"mysql" reserved for future use.
	DBDSN    string // sqlite file path, or connection string for other drivers.

	Retry RetryConfig
}

// RetryConfig controls the retry/backoff behavior applied to database operations.
type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

// Load reads configuration from environment variables, applying sane defaults
// for local development so the server can run with zero setup.
func Load() (Config, error) {
	cfg := Config{
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		CORSAllowOrigin: getEnv("CORS_ALLOW_ORIGIN", "http://localhost:3000"),
		DBDriver:        getEnv("DB_DRIVER", "sqlite"),
		DBDSN:           getEnv("DB_DSN", "file:data/stock-stonks.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"),
		Retry: RetryConfig{
			MaxAttempts: getEnvInt("DB_RETRY_MAX_ATTEMPTS", 5),
			BaseDelay:   getEnvDuration("DB_RETRY_BASE_DELAY", 50*time.Millisecond),
			MaxDelay:    getEnvDuration("DB_RETRY_MAX_DELAY", 2*time.Second),
		},
	}

	if cfg.Retry.MaxAttempts < 1 {
		return Config{}, fmt.Errorf("config: DB_RETRY_MAX_ATTEMPTS must be >= 1, got %d", cfg.Retry.MaxAttempts)
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}
