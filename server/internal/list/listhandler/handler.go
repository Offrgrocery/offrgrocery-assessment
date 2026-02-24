package listhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"offgrocery-assessment/internal/httputil"
	"offgrocery-assessment/internal/list/listservice"
)

type Handler interface {
	Routes() chi.Router
	CreateList(w http.ResponseWriter, r *http.Request)
	GetLists(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)
	DeleteList(w http.ResponseWriter, r *http.Request)
	AddItems(w http.ResponseWriter, r *http.Request)
	RemoveItems(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service listservice.Service
}

func New(service listservice.Service) *handler {
	return &handler{service: service}
}

func (h *handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateList)
	r.Get("/", h.GetLists)
	r.Get("/{id}", h.GetList)
	r.Delete("/{id}", h.DeleteList)
	r.Post("/{id}/items", h.AddItems)
	r.Delete("/{id}/items", h.RemoveItems)
	return r
}

type createListRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type addItemsRequest struct {
	ItemIDs []int `json:"item_ids"`
}

type removeItemsRequest struct {
	ItemIDs []int `json:"item_ids"`
}

func (h *handler) CreateList(w http.ResponseWriter, r *http.Request) {
	var req createListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.UserID == 0 || req.Name == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "user_id and name are required"})
		return
	}

	list, err := h.service.CreateList(r.Context(), req.UserID, req.Name)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrorResponse{Error: "failed to create list"})
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, list)
}

func (h *handler) GetLists(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "user_id query param is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid user_id"})
		return
	}

	lists, err := h.service.GetListsByUserID(r.Context(), userID)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrorResponse{Error: "failed to fetch lists"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, lists)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid list id"})
		return
	}

	list, err := h.service.GetListByID(r.Context(), id)
	if err != nil {
		httputil.WriteJSON(w, http.StatusNotFound, httputil.ErrorResponse{Error: "list not found"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *handler) DeleteList(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid list id"})
		return
	}

	if err := h.service.DeleteList(r.Context(), id); err != nil {
		httputil.WriteJSON(w, http.StatusNotFound, httputil.ErrorResponse{Error: "list not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) AddItems(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	listID, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid list id"})
		return
	}

	var req addItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid request body"})
		return
	}

	if len(req.ItemIDs) == 0 {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "item_ids are required"})
		return
	}

	list, err := h.service.AddItemsToList(r.Context(), listID, req.ItemIDs)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrorResponse{Error: "failed to add items to list"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, list)
}

func (h *handler) RemoveItems(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	listID, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid list id"})
		return
	}

	var req removeItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid request body"})
		return
	}

	if len(req.ItemIDs) == 0 {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "item_ids are required"})
		return
	}

	list, err := h.service.RemoveItemsFromList(r.Context(), listID, req.ItemIDs)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrorResponse{Error: "failed to remove items from list"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, list)
}
