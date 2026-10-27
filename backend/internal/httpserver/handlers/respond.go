package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
)

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		slog.Error("httpserver: encode response", "error", err)
	}
}

type errorBody struct {
	Error string `json:"error"`
}

var errInvalidIDParam = errors.New("invalid id parameter")

// writeError maps a domain/repository error to the right HTTP status instead
// of leaking every failure as a 500, so API clients can branch on status
// code alone.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		writeJSON(w, http.StatusNotFound, errorBody{Error: "not found"})
	case errors.Is(err, repository.ErrDuplicateTicker):
		writeJSON(w, http.StatusConflict, errorBody{Error: err.Error()})
	case isValidationError(err):
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
	default:
		slog.Error("httpserver: internal error", "error", err)
		writeJSON(w, http.StatusInternalServerError, errorBody{Error: "internal server error"})
	}
}

func isValidationError(err error) bool {
	switch {
	case errors.Is(err, domain.ErrInvalidTicker),
		errors.Is(err, domain.ErrInvalidName),
		errors.Is(err, domain.ErrInvalidType),
		errors.Is(err, domain.ErrInvalidAssetID),
		errors.Is(err, domain.ErrInvalidOperation),
		errors.Is(err, domain.ErrInvalidQuantity),
		errors.Is(err, domain.ErrInvalidUnitPrice),
		errors.Is(err, domain.ErrNegativeFees),
		errors.Is(err, domain.ErrInvalidDate):
		return true
	default:
		return false
	}
}
