package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nataalka/splitni-to/internal/domain"
	"github.com/nataalka/splitni-to/internal/domain/models"
	"github.com/nataalka/splitni-to/internal/domain/services"
)

type FriendHandler struct {
	friendService *services.FriendService
}

func NewFriendHandler(fs *services.FriendService) *FriendHandler {
	return &FriendHandler{friendService: fs}
}

// AddFriend godoc
// @Summary      Send friend request
// @Description  Send a new friend request using the target user's email
// @Tags         friends
// @Security     BearerAuth
// @Accept       json
// @Param        request body models.FriendRequest true "Friend's email"
// @Success      201 "Created"
// @Failure      404 {object} domain.AppError "User not found"
// @Router       /friends [post]
func (h *FriendHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	var req models.FriendRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	err = h.friendService.SendFriendRequest(ctx, userID, req.Email)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// AcceptFriend godoc
// @Summary      Accept friend request
// @Description  Accept a pending friend request from another user
// @Tags         friends
// @Security     BearerAuth
// @Param        id path string true "Requester User ID"
// @Success      200 "OK"
// @Router       /friends/{id}/accept [put]
func (h *FriendHandler) AcceptFriend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	requesterID, err := parseUUID(r, "id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	err = h.friendService.AcceptFriendRequest(ctx, userID, requesterID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteFriendship godoc
// @Summary      Remove friend
// @Description  Remove a user from friends list or decline a request
// @Tags         friends
// @Security     BearerAuth
// @Param        id path string true "Friend User ID"
// @Success      204 "No Content"
// @Router       /friends/{id} [delete]
func (h *FriendHandler) DeleteFriendship(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	friendID, err := parseUUID(r, "id")
	if err != nil {
		respondWithError(w, err)
	}

	err = h.friendService.DeleteFriendship(ctx, userID, friendID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPendingRequests godoc
// @Summary      List pending requests
// @Description  Get all incoming friend requests that are waiting for approval
// @Tags         friends
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} models.User
// @Router       /friends/pending [get]
func (h *FriendHandler) ListPendingRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	pending, err := h.friendService.GetPendingRequests(ctx, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, pending)
}

// ListFriends godoc
// @Summary      List all friends
// @Description  Get a list of all accepted friends for the authenticated user
// @Tags         friends
// @Security     BearerAuth
// @Produce      json
// @Success      200 {array} models.User
// @Router       /friends [get]
func (h *FriendHandler) ListFriends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	friends, err := h.friendService.GetFriendsList(ctx, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, friends)
}
