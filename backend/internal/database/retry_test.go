package database_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/config"
	"github.com/danillomeira/stock-stonks/backend/internal/database"
)

func TestIsRetryable(t *testing.T) {
	cases := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{errors.New("database is locked"), true},
		{errors.New("SQLITE_BUSY: database table is locked"), true},
		{errors.New("read: connection reset by peer"), true},
		{errors.New("write: broken pipe"), true},
		{errors.New("UNIQUE constraint failed: assets.ticker"), false},
		{errors.New("record not found"), false},
	}

	for _, tc := range cases {
		got := database.IsRetryable(tc.err)
		if got != tc.want {
			t.Errorf("IsRetryable(%v) = %v, want %v", tc.err, got, tc.want)
		}
	}
}

func TestWithRetry_SucceedsAfterTransientErrors(t *testing.T) {
	cfg := config.RetryConfig{MaxAttempts: 5, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}

	attempts := 0
	err := database.WithRetry(context.Background(), cfg, func() error {
		attempts++
		if attempts < 3 {
			return errors.New("database is locked")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("WithRetry() error = %v, want nil", err)
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
}

func TestWithRetry_DoesNotRetryPermanentErrors(t *testing.T) {
	cfg := config.RetryConfig{MaxAttempts: 5, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}
	permanent := errors.New("UNIQUE constraint failed")

	attempts := 0
	err := database.WithRetry(context.Background(), cfg, func() error {
		attempts++
		return permanent
	})

	if !errors.Is(err, permanent) {
		t.Fatalf("WithRetry() error = %v, want %v", err, permanent)
	}
	if attempts != 1 {
		t.Fatalf("attempts = %d, want 1 (no retry on permanent error)", attempts)
	}
}

func TestWithRetry_ExhaustsAttempts(t *testing.T) {
	cfg := config.RetryConfig{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 10 * time.Millisecond}

	attempts := 0
	err := database.WithRetry(context.Background(), cfg, func() error {
		attempts++
		return errors.New("database is locked")
	})

	if !errors.Is(err, database.ErrRetryExhausted) {
		t.Fatalf("WithRetry() error = %v, want wrapped ErrRetryExhausted", err)
	}
	if attempts != cfg.MaxAttempts {
		t.Fatalf("attempts = %d, want %d", attempts, cfg.MaxAttempts)
	}
}

func TestWithRetry_RespectsContextCancellation(t *testing.T) {
	cfg := config.RetryConfig{MaxAttempts: 10, BaseDelay: 50 * time.Millisecond, MaxDelay: time.Second}
	ctx, cancel := context.WithCancel(context.Background())

	attempts := 0
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()

	err := database.WithRetry(ctx, cfg, func() error {
		attempts++
		return errors.New("database is locked")
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("WithRetry() error = %v, want context.Canceled", err)
	}
	if attempts >= cfg.MaxAttempts {
		t.Fatalf("attempts = %d, expected cancellation to stop retries before exhausting attempts", attempts)
	}
}
