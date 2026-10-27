package database

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/config"
)

// ErrRetryExhausted is returned when every retry attempt failed.
var ErrRetryExhausted = errors.New("database: retry attempts exhausted")

// IsRetryable decides whether an error from the database driver represents a
// transient condition worth retrying (e.g. SQLite reporting the database is
// locked under concurrent writers, or a dropped connection to a networked
// database). Anything else - constraint violations, not-found, bad input -
// is permanent and must fail immediately.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "database is locked"):
		return true
	case strings.Contains(msg, "database table is locked"):
		return true
	case strings.Contains(msg, "busy"):
		return true
	case strings.Contains(msg, "connection reset"):
		return true
	case strings.Contains(msg, "broken pipe"):
		return true
	case strings.Contains(msg, "too many connections"):
		return true
	default:
		return false
	}
}

// WithRetry executes fn, retrying with exponential backoff and jitter when
// fn's error is classified as retryable by IsRetryable. It exists so every
// repository implementation shares one retry policy instead of each one
// reimplementing backoff math - and so swapping SQLite for Postgres/MySQL
// later only means teaching IsRetryable their driver-specific transient
// errors, not rewriting every query call site.
func WithRetry(ctx context.Context, cfg config.RetryConfig, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		if !IsRetryable(lastErr) {
			return lastErr
		}

		delay := backoffDelay(attempt, cfg.BaseDelay, cfg.MaxDelay)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("%w after %d attempts: %v", ErrRetryExhausted, cfg.MaxAttempts, lastErr)
}

// backoffDelay computes an exponential delay with +/-25% jitter, capped at maxDelay.
func backoffDelay(attempt int, base, maxDelay time.Duration) time.Duration {
	raw := float64(base) * math.Pow(2, float64(attempt))
	if raw > float64(maxDelay) {
		raw = float64(maxDelay)
	}
	jitter := raw * (0.75 + rand.Float64()*0.5) // 75%-125% of raw
	return time.Duration(jitter)
}
