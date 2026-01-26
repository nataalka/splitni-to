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

type FriendHandler struct {
	friendService *services.FriendService
}

func NewFriendHandler(fs *services.FriendService) *FriendHandler {
	return &FriendHandler{friendService: fs}
}

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

func (h *FriendHandler) AcceptFriend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	requesterIDStr := chi.URLParam(r, "id")
	requesterID, err := uuid.Parse(requesterIDStr)
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid user id", err))
		return
	}

	err = h.friendService.AcceptFriendRequest(ctx, userID, requesterID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *FriendHandler) DeleteFriendship(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, err)
		return
	}

	friendIDStr := chi.URLParam(r, "id")
	friendID, _ := uuid.Parse(friendIDStr)

	err = h.friendService.DeleteFriendship(ctx, userID, friendID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

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
