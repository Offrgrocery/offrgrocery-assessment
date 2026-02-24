package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"offgrocery-assessment/internal/auth/authservice"
	"offgrocery-assessment/internal/httputil"
)

type Handler interface {
	Routes() chi.Router
	CreateUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	service authservice.Service
}

func New(service authservice.Service) *handler {
	return &handler{service: service}
}

func (h *handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateUser)
	return r
}

type createUserRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type createUserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "invalid request body"})
		return
	}

	if req.Email == "" || req.Name == "" {
		httputil.WriteJSON(w, http.StatusBadRequest, httputil.ErrorResponse{Error: "email and name are required"})
		return
	}

	id, err := h.service.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, httputil.ErrorResponse{Error: "failed to create user"})
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, createUserResponse{
		ID:    id,
		Email: req.Email,
		Name:  req.Name,
	})
}
