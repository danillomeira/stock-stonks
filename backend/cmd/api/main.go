// Command api runs the Stock-Stonks HTTP API.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	assetRepo := sqlite.NewAssetRepository(db, cfg.Retry)
	operationRepo := sqlite.NewOperationRepository(db, cfg.Retry)

	svc := httpserver.Services{
		Assets:     service.NewAssetService(assetRepo),
		Operations: service.NewOperationService(operationRepo, assetRepo),
		Positions:  service.NewPositionService(operationRepo),
	}

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      httpserver.NewRouter(svc, cfg.CORSAllowOrigin),
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
