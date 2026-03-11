package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/backend/internal/domain"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
	repositories2 "github.com/nataalka/splitni-to/backend/internal/domain/repositories"
)

type FriendService struct {
	friendRepo repositories2.FriendRepository
	userRepo   repositories2.UserRepository
}

func NewFriendService(fr repositories2.FriendRepository, ur repositories2.UserRepository) *FriendService {
	return &FriendService{
		friendRepo: fr,
		userRepo:   ur,
	}
}

func (s *FriendService) SendFriendRequest(ctx context.Context, senderID uuid.UUID, friendEmail string) error {
	friend, err := s.userRepo.GetByEmail(ctx, friendEmail)
	if err != nil {
		return domain.NewNotFoundError("user with this email not found")
	}

	if senderID == friend.ID {
		return domain.NewValidationError("cannot add yourself as a friend")
	}

	return s.friendRepo.CreateFriendRequest(ctx, senderID, friend.ID)
}

func (s *FriendService) AcceptFriendRequest(ctx context.Context, userID, requesterID uuid.UUID) error {
	return s.friendRepo.AcceptFriendRequest(ctx, userID, requesterID)
}

func (s *FriendService) DeleteFriendship(ctx context.Context, id1, id2 uuid.UUID) error {
	return s.friendRepo.DeleteFriendship(ctx, id1, id2)
}

func (s *FriendService) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]models.User, error) {
	return s.friendRepo.GetPendingRequests(ctx, userID)
}

func (s *FriendService) GetFriendsList(ctx context.Context, userID uuid.UUID) ([]models.User, error) {
	return s.friendRepo.GetFriends(ctx, userID)
}
