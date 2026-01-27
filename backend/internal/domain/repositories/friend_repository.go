package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
)

type FriendRepository interface {
	CreateFriendRequest(ctx context.Context, id1, id2 uuid.UUID) error
	AcceptFriendRequest(ctx context.Context, id1, id2 uuid.UUID) error
	DeleteFriendship(ctx context.Context, id1, id2 uuid.UUID) error
	GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]models.User, error)
	GetFriends(ctx context.Context, userID uuid.UUID) ([]models.User, error)
}
