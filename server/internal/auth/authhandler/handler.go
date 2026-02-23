package authhandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"offgrocery-assessment/internal/auth/authservice"
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

type errorResponse struct {
	Error string `json:"error"`
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request body"})
		return
	}

	if req.Email == "" || req.Name == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "email and name are required"})
		return
	}

	id, err := h.service.CreateUser(r.Context(), req.Email, req.Name)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse{Error: "failed to create user"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createUserResponse{
		ID:    id,
		Email: req.Email,
		Name:  req.Name,
	})
}
