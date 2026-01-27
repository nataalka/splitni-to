package repositories

import (
	"context"

	"github.com/google/uuid"
	models2 "github.com/nataalka/splitni-to/backend/internal/domain/models"
)

type GroupRepository interface {
	Create(ctx context.Context, group *models2.Group, userID uuid.UUID) error
	AddMember(ctx context.Context, groupID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, groupID, userID uuid.UUID) error
	IsMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error)
	GetByID(ctx context.Context, groupID uuid.UUID) (*models2.Group, error)
	GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]models2.User, error)
	GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]models2.Group, error)
}
