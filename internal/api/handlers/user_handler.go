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
}

func NewUserHandler(us *services.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateUser
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, domain.NewValidationError("invalid json format", err))
		return
	}

	user, err := h.userService.Register(ctx, req)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		RespondWithError(w, domain.NewValidationError("invalid uuid format", err))
		return
	}

	user, err := h.userService.GetByID(r.Context(), id)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, user)
}
