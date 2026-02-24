package itemhandler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"offgrocery-assessment/internal/httputil"
	"offgrocery-assessment/internal/item/itemservice"
)

type Handler interface {
	Routes() chi.Router
	GetItem(w http.ResponseWriter, r *http.Request)
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
	r.Get("/{id}", h.GetItem)
	return r
}

func (h *handler) GetItem(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid item id"})
		return
	}

	item, err := h.service.GetItemByID(r.Context(), id)
	if err != nil {
		httputil.WriteJSON(w, http.StatusNotFound, httputil.ErrorResponse{Error: "item not found"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, item)
}

func (h *handler) SearchWithLimit(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "q query param is required"})
		return
	}

	limit := 20
	if countStr := r.URL.Query().Get("count"); countStr != "" {
		parsed, err := strconv.Atoi(countStr)
		if err != nil || parsed < 1 {
			httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid count"})
			return
		}
		limit = parsed
	}

	items, err := h.service.SearchWithLimit(r.Context(), query, limit)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrorResponse{Error: "failed to search items"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, items)
}
