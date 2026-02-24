package storehandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"offgrocery-assessment/internal/store/storeservice"
)

type Handler interface {
	Routes() chi.Router
	GetStore(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service storeservice.Service
}

func New(service storeservice.Service) *handler {
	return &handler{service: service}
}

func (h *handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/{id}", h.GetStore)
	return r
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *handler) GetStore(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid store id"})
		return
	}

	store, err := h.service.GetStoreByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "store not found"})
		return
	}

	writeJSON(w, http.StatusOK, store)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
