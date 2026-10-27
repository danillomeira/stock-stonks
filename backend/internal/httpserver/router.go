// Package httpserver wires domain services to HTTP routes.
package httpserver

import (
	"net/http"

	"github.com/danillomeira/stock-stonks/backend/internal/auth"
	"github.com/danillomeira/stock-stonks/backend/internal/httpserver/handlers"
	"github.com/danillomeira/stock-stonks/backend/internal/httpserver/middleware"
	"github.com/danillomeira/stock-stonks/backend/internal/service"
)

// Services groups every service the router needs to build handlers.
type Services struct {
	Auth       *service.AuthService
	Assets     *service.AssetService
	Operations *service.OperationService
	Positions  *service.PositionService
}

// RouterConfig carries the non-service settings the router needs.
type RouterConfig struct {
	CORSAllowOrigin string
	CookieSecure    bool
	Sessions        *auth.SessionManager
}

// NewRouter builds the full HTTP handler tree, including middleware.
func NewRouter(svc Services, cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	authHandler := handlers.NewAuthHandler(svc.Auth, cfg.CORSAllowOrigin, cfg.CookieSecure)
	assetHandler := handlers.NewAssetHandler(svc.Assets)
	operationHandler := handlers.NewOperationHandler(svc.Operations)
	positionHandler := handlers.NewPositionHandler(svc.Positions)

	requireAuth := middleware.RequireAuth(cfg.Sessions)

	mux.HandleFunc("GET /health", handlers.Health)

	// Public: starts/completes the OAuth flow, so by definition runs before a session exists.
	mux.HandleFunc("GET /api/v1/auth/google/login", authHandler.Login)
	mux.HandleFunc("GET /api/v1/auth/google/callback", authHandler.Callback)
	mux.HandleFunc("POST /api/v1/auth/logout", authHandler.Logout)

	// Protected: every data route, plus "who am I" which needs a verified session.
	mux.Handle("GET /api/v1/auth/me", requireAuth(http.HandlerFunc(authHandler.Me)))

	mux.Handle("POST /api/v1/assets", requireAuth(http.HandlerFunc(assetHandler.Create)))
	mux.Handle("GET /api/v1/assets", requireAuth(http.HandlerFunc(assetHandler.List)))
	mux.Handle("GET /api/v1/assets/{id}", requireAuth(http.HandlerFunc(assetHandler.Get)))
	mux.Handle("PUT /api/v1/assets/{id}", requireAuth(http.HandlerFunc(assetHandler.Update)))
	mux.Handle("DELETE /api/v1/assets/{id}", requireAuth(http.HandlerFunc(assetHandler.Delete)))

	mux.Handle("POST /api/v1/operations", requireAuth(http.HandlerFunc(operationHandler.Create)))
	mux.Handle("GET /api/v1/operations", requireAuth(http.HandlerFunc(operationHandler.List)))
	mux.Handle("GET /api/v1/operations/{id}", requireAuth(http.HandlerFunc(operationHandler.Get)))
	mux.Handle("PUT /api/v1/operations/{id}", requireAuth(http.HandlerFunc(operationHandler.Update)))
	mux.Handle("DELETE /api/v1/operations/{id}", requireAuth(http.HandlerFunc(operationHandler.Delete)))

	mux.Handle("GET /api/v1/positions", requireAuth(http.HandlerFunc(positionHandler.List)))
	mux.Handle("GET /api/v1/positions/{id}", requireAuth(http.HandlerFunc(positionHandler.Get)))

	return middleware.Logging(middleware.CORS(cfg.CORSAllowOrigin)(mux))
}
