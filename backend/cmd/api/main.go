// Command api runs the Stock-Stonks HTTP API.
package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/auth"
	"github.com/danillomeira/stock-stonks/backend/internal/config"
	"github.com/danillomeira/stock-stonks/backend/internal/database"
	"github.com/danillomeira/stock-stonks/backend/internal/httpserver"
	"github.com/danillomeira/stock-stonks/backend/internal/repository/sqlite"
	"github.com/danillomeira/stock-stonks/backend/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	db, err := database.Open(cfg)
	if err != nil {
		return err
	}
	if err := sqlite.AutoMigrate(db); err != nil {
		return err
	}

	userRepo := sqlite.NewUserRepository(db, cfg.Retry)
	assetRepo := sqlite.NewAssetRepository(db, cfg.Retry)
	operationRepo := sqlite.NewOperationRepository(db, cfg.Retry)

	sessionSecret := cfg.Auth.SessionSecret
	if sessionSecret == "" {
		sessionSecret, err = generateEphemeralSecret()
		if err != nil {
			return err
		}
		slog.Warn("SESSION_SECRET not set - using a random secret for this process; existing sessions will be invalidated on every restart")
	}
	sessions := auth.NewSessionManager(sessionSecret, cfg.Auth.SessionTTL)

	googleOAuth, err := auth.NewGoogleOAuth(ctx, cfg.Auth.GoogleClientID, cfg.Auth.GoogleClientSecret, cfg.Auth.GoogleRedirectURL)
	if err != nil {
		return err
	}

	svc := httpserver.Services{
		Auth:       service.NewAuthService(googleOAuth, sessions, userRepo),
		Assets:     service.NewAssetService(assetRepo),
		Operations: service.NewOperationService(operationRepo, assetRepo),
		Positions:  service.NewPositionService(operationRepo),
	}

	router := httpserver.NewRouter(svc, httpserver.RouterConfig{
		CORSAllowOrigin: cfg.CORSAllowOrigin,
		CookieSecure:    cfg.Auth.CookieSecure,
		Sessions:        sessions,
	})

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("http server listening", "port", cfg.HTTPPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	stopCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return err
	case <-stopCtx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}

func generateEphemeralSecret() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
