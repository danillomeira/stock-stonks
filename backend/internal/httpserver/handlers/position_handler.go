package handlers

import (
	"net/http"

	"github.com/danillomeira/stock-stonks/backend/internal/service"
)

type PositionHandler struct {
	positions *service.PositionService
}

func NewPositionHandler(positions *service.PositionService) *PositionHandler {
	return &PositionHandler{positions: positions}
}

func (h *PositionHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	positions, err := h.positions.All(r.Context(), userID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, positions)
}

func (h *PositionHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	position, err := h.positions.ForAsset(r.Context(), userID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, position)
}
