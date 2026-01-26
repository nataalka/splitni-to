package handlers

import (
	"encoding/json"
	"net/http"

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

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	var req models.CreateGroupRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	group, err := h.groupService.CreateGroup(ctx, req.Name, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, group)
}

func (h *GroupHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groupID, err := parseUUID(r, "id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	var req models.AddMemberRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	err = h.groupService.AddMemberToGroup(ctx, userID, groupID, req.UserID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groupID, err := parseUUID(r, "id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	targetID, err := parseUUID(r, "userID")
	if err != nil {
		respondWithError(w, err)
		return
	}

	err = h.groupService.RemoveMember(ctx, userID, groupID, targetID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groupID, err := parseUUID(r, "id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	group, err := h.groupService.GetByID(ctx, userID, groupID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, group)
}

func (h *GroupHandler) ListUsersGroups(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groups, err := h.groupService.GetUserGroups(ctx, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, groups)
}

func (h *GroupHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groupID, err := parseUUID(r, "id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	members, err := h.groupService.GetMembers(ctx, userID, groupID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, members)
}
