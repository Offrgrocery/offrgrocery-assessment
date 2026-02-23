package itemhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"offgrocery-assessment/internal/item/itemservice"

	"github.com/go-chi/chi/v5"
)

type Handler interface {
	Routes() chi.Router
	SearchWithLimit(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service itemservice.Service
}

func New(service itemservice.Service) *handler {
	return &handler{service: service}
}

func (h *handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/search", h.SearchWithLimit)
	return r
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *handler) SearchWithLimit(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "q query param is required"})
		return
	}

	limit := 20
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		parsed, err := strconv.Atoi(countStr)
		if err != nil || parsed < 1 {
			writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid count"})
			return
		}
		limit = parsed
	}

	items, err := h.service.SearchWithLimit(r.Context(), query, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to search items"})
		return
	}

	writeJSON(w, http.StatusOK, items)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
