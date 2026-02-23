package listhandler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"offgrocery-assessment/internal/list/listservice"
)

type Handler interface {
	Routes() chi.Router
	CreateList(w http.ResponseWriter, r *http.Request)
	GetLists(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)
	AddItems(w http.ResponseWriter, r *http.Request)
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
	r.Post("/{id}/items", h.AddItems)
	return r
}

type createListRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
}

type addItemsRequest struct {
	ItemIDs []int `json:"item_ids"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (h *handler) CreateList(w http.ResponseWriter, r *http.Request) {
	var req createListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	if req.UserID == 0 || req.Name == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "user_id and name are required"})
		return
	}

	list, err := h.service.CreateList(r.Context(), req.UserID, req.Name)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to create list"})
		return
	}

	writeJSON(w, http.StatusCreated, list)
}

func (h *handler) GetLists(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "user_id query param is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid user_id"})
		return
	}

	lists, err := h.service.GetListsByUserID(r.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to fetch lists"})
		return
	}

	writeJSON(w, http.StatusOK, lists)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid list id"})
		return
	}

	list, err := h.service.GetListByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "list not found"})
		return
	}

	writeJSON(w, http.StatusOK, list)
}

func (h *handler) AddItems(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	listID, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid list id"})
		return
	}

	var req addItemsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}

	if len(req.ItemIDs) == 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "item_ids are required"})
		return
	}

	list, err := h.service.AddItemsToList(r.Context(), listID, req.ItemIDs)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "failed to add items to list"})
		return
	}

	writeJSON(w, http.StatusOK, list)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
