package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nataalka/splitni-to/internal/domain"
	"github.com/nataalka/splitni-to/internal/domain/models"
)

type GroupRepository struct {
	db *sqlx.DB
}

func NewPostgresGroupRepository(db *sqlx.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *models.Group, userID uuid.UUID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return domain.NewInternalError("failed to begin transaction", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "INSERT INTO groups (id, name, created_at) VALUES ($1, $2, $3)", group.ID, group.Name, group.CreatedAt)
	if err != nil {
		return domain.NewInternalError("could not create group", err)
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO group_members (group_id, user_id) VALUES ($1, $2)", group.ID, userID)
	if err != nil {
		return domain.NewInternalError("could not insert group member", err)
	}

	return tx.Commit()
}
