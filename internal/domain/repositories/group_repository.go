package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain/models"
)

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group, userID uuid.UUID) error
	AddMember(ctx context.Context, groupID, userID uuid.UUID) error
	RemoveMember(ctx context.Context, groupID, userID uuid.UUID) error
	IsMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error)
	GetByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error)
	GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]models.User, error)
	GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Group, error)
}
