package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}
