// Package httpserver wires domain services to HTTP routes.
package httpserver

import (
	"net/http"

	"github.com/danillomeira/stock-stonks/backend/internal/httpserver/handlers"
	"github.com/danillomeira/stock-stonks/backend/internal/httpserver/middleware"
	"github.com/danillomeira/stock-stonks/backend/internal/service"
)

// Services groups every service the router needs to build handlers.
type Services struct {
	Assets     *service.AssetService
	Operations *service.OperationService
	Positions  *service.PositionService
}

// NewRouter builds the full HTTP handler tree, including middleware.
func NewRouter(svc Services, corsAllowOrigin string) http.Handler {
	mux := http.NewServeMux()

	assetHandler := handlers.NewAssetHandler(svc.Assets)
	operationHandler := handlers.NewOperationHandler(svc.Operations)
	positionHandler := handlers.NewPositionHandler(svc.Positions)

	mux.HandleFunc("GET /health", handlers.Health)

	mux.HandleFunc("POST /api/v1/assets", assetHandler.Create)
	mux.HandleFunc("GET /api/v1/assets", assetHandler.List)
	mux.HandleFunc("GET /api/v1/assets/{id}", assetHandler.Get)
	mux.HandleFunc("PUT /api/v1/assets/{id}", assetHandler.Update)
	mux.HandleFunc("DELETE /api/v1/assets/{id}", assetHandler.Delete)

	mux.HandleFunc("POST /api/v1/operations", operationHandler.Create)
	mux.HandleFunc("GET /api/v1/operations", operationHandler.List)
	mux.HandleFunc("GET /api/v1/operations/{id}", operationHandler.Get)
	mux.HandleFunc("PUT /api/v1/operations/{id}", operationHandler.Update)
	mux.HandleFunc("DELETE /api/v1/operations/{id}", operationHandler.Delete)

	mux.HandleFunc("GET /api/v1/positions", positionHandler.List)
	mux.HandleFunc("GET /api/v1/positions/{id}", positionHandler.Get)

	return middleware.Logging(middleware.CORS(corsAllowOrigin)(mux))
}
