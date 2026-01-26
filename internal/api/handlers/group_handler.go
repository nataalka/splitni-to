package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain"
	"github.com/nataalka/splitni-to/internal/domain/models"
	"github.com/nataalka/splitni-to/internal/domain/services"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler(gs *services.GroupService) *GroupHandler {
	return &GroupHandler{groupService: gs}
}

func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	_, claims, _ := jwtauth.FromContext(r.Context())
	sub, ok := claims["sub"].(string)
	if !ok {
		RespondWithError(w, domain.NewAuthError("invalid token claims"))
		return
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		RespondWithError(w, domain.NewInternalError("invalid user id in token", err))
		return
	}

	var req models.CreateGroupRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		RespondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	group, err := h.groupService.CreateGroup(ctx, req.Name, userID)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	WriteJSON(w, http.StatusCreated, group)
}
