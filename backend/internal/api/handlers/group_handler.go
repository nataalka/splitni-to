package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nataalka/splitni-to/backend/internal/domain"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
	"github.com/nataalka/splitni-to/backend/internal/domain/services"
)

type GroupHandler struct {
	groupService *services.GroupService
}

func NewGroupHandler(gs *services.GroupService) *GroupHandler {
	return &GroupHandler{groupService: gs}
}

// Create godoc
// @Summary      Create a new group
// @Description  Create a group and automatically add the creator as a member
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        group body models.CreateGroupRequest true "Group name"
// @Success      201 {object} models.Group
// @Router       /groups [post]
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

	group, err := h.groupService.CreateGroup(ctx, req, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, group)
}

// Update godoc
// @Summary      Update group details
// @Tags         groups
// @Param        group_id path string true "Group ID"
// @Param        group body models.CreateGroupRequest true "New details"
// @Router       /groups/{group_id} [patch]
func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := getUserIDFromContext(ctx)

	groupID, err := parseUUID(r, "group_id")
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

	err = h.groupService.UpdateGroup(ctx, userID, groupID, req)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary      Delete a group
// @Tags         groups
// @Param        group_id path string true "Group ID"
// @Router       /groups/{group_id} [delete]
func (h *GroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := getUserIDFromContext(ctx)

	groupID, err := parseUUID(r, "group_id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	err = h.groupService.DeleteGroup(ctx, userID, groupID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddMember godoc
// @Summary      Add a member to a group
// @Description  Add a new user to an existing group by their User ID
// @Tags         groups
// @Security     BearerAuth
// @Accept       json
// @Param        group_id path string true "Group ID"
// @Param        request body models.AddMemberRequest true "User ID to add"
// @Success      204 "No Content"
// @Failure      403 {object} domain.AppError "Permission denied"
// @Failure      404 {object} domain.AppError "Group not found"
// @Router       /groups/{group_id}/members [post]
func (h *GroupHandler) AddMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groupID, err := parseUUID(r, "group_id")
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

// RemoveMember godoc
// @Summary      Remove a member from a group
// @Description  Remove a user from the group. Can be used for leaving or kicking members.
// @Tags         groups
// @Security     BearerAuth
// @Param        group_id path string true "Group ID"
// @Param        user_id path string true "User ID to remove"
// @Success      204 "No Content"
// @Failure      403 {object} domain.AppError
// @Router       /groups/{group_id}/members/{user_id} [delete]
func (h *GroupHandler) RemoveMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groupID, err := parseUUID(r, "group_id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	targetID, err := parseUUID(r, "user_id")
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

// GetByID godoc
// @Summary      Get group details
// @Description  Fetch details of a specific group if the user is a member
// @Tags         groups
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Group ID"
// @Success      200 {object} models.Group
// @Failure      403 {object} domain.AppError
// @Router       /groups/{group_id} [get]
func (h *GroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groupID, err := parseUUID(r, "group_id")
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

// ListUsersGroups godoc
// @Summary      List my groups
// @Description  Get a list of all groups the authenticated user belongs to
// @Tags         groups
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} models.Group
// @Router       /groups [get]
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

// ListMembers godoc
// @Summary      List group members
// @Description  Get a list of all users in a specific group
// @Tags         groups
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "Group ID"
// @Success      200 {array} models.User
// @Router       /groups/{group_id}/members [get]
func (h *GroupHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	groupID, err := parseUUID(r, "group_id")
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
