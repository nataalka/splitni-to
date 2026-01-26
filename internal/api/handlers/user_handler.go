package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain"
	"github.com/nataalka/splitni-to/internal/domain/models"
	"github.com/nataalka/splitni-to/internal/domain/services"
)

type UserHandler struct {
	userService *services.UserService
	jwtSecret   string
}

func NewUserHandler(us *services.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		userService: us,
		jwtSecret:   jwtSecret,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	user, err := h.userService.Register(ctx, req)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	if req.Email == "" || req.Password == "" {
		RespondWithError(w, domain.NewValidationError("email and password are required"))
		return
	}

	res, err := h.userService.Login(ctx, req.Email, req.Password, h.jwtSecret)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, res)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, domain.NewValidationError("invalid uuid format", err))
		return
	}

	user, err := h.userService.GetByID(ctx, id)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, user)
}
