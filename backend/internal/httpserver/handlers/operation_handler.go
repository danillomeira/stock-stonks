package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
	"github.com/danillomeira/stock-stonks/backend/internal/service"
	"github.com/shopspring/decimal"
)

type OperationHandler struct {
	operations *service.OperationService
}

func NewOperationHandler(operations *service.OperationService) *OperationHandler {
	return &OperationHandler{operations: operations}
}

type operationRequest struct {
	AssetID   uint                 `json:"asset_id"`
	Date      time.Time            `json:"date"`
	Type      domain.OperationType `json:"type"`
	Quantity  decimal.Decimal      `json:"quantity"`
	UnitPrice decimal.Decimal      `json:"unit_price"`
	Fees      decimal.Decimal      `json:"fees"`
	Broker    string               `json:"broker"`
}

func (req operationRequest) toDomain() *domain.Operation {
	return &domain.Operation{
		AssetID:   req.AssetID,
		Date:      req.Date,
		Type:      req.Type,
		Quantity:  req.Quantity,
		UnitPrice: req.UnitPrice,
		Fees:      req.Fees,
		Broker:    req.Broker,
	}
}

func (h *OperationHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "invalid request body"})
		return
	}

	op := req.toDomain()
	if err := h.operations.Create(r.Context(), userID, op); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, op)
}

func (h *OperationHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	op, err := h.operations.Get(r.Context(), userID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, op)
}

func (h *OperationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	filter := repository.OperationFilter{UserID: userID}
	if v := r.URL.Query().Get("asset_id"); v != "" {
		id, err := parseUintQuery(v)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorBody{Error: "invalid asset_id query param"})
			return
		}
		filter.AssetID = &id
	}
	if v := r.URL.Query().Get("from"); v != "" {
		t, err := time.Parse(time.DateOnly, v)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorBody{Error: "invalid from query param, expected YYYY-MM-DD"})
			return
		}
		filter.From = &t
	}
	if v := r.URL.Query().Get("to"); v != "" {
		t, err := time.Parse(time.DateOnly, v)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errorBody{Error: "invalid to query param, expected YYYY-MM-DD"})
			return
		}
		filter.To = &t
	}

	ops, err := h.operations.List(r.Context(), filter)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, ops)
}

func (h *OperationHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	var req operationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: "invalid request body"})
		return
	}

	op := req.toDomain()
	op.ID = id
	if err := h.operations.Update(r.Context(), userID, op); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, op)
}

func (h *OperationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := requireUserID(w, r)
	if !ok {
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorBody{Error: err.Error()})
		return
	}

	if err := h.operations.Delete(r.Context(), userID, id); err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusNoContent, nil)
}

func parseUintQuery(v string) (uint, error) {
	id, err := strconv.ParseUint(v, 10, 64)
	return uint(id), err
}
