package storehandler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"offgrocery-assessment/internal/httputil"
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

func (h *handler) GetStore(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid store id"})
		return
	}

	store, err := h.service.GetStoreByID(r.Context(), id)
	if err != nil {
		httputil.WriteJSON(w, http.StatusNotFound, httputil.ErrorResponse{Error: "store not found"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, store)
}
