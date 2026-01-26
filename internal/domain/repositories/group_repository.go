package repositories

import (
	"context"

	"github.com/nataalka/splitni-to/internal/domain/models"
)

type GroupRepository interface {
	Create(ctx context.Context, group *models.Group) error
}
