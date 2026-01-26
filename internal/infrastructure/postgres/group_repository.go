package postgres

import (
	"context"
	"database/sql"
	"errors"

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

func (r *GroupRepository) AddMember(ctx context.Context, groupID, userID uuid.UUID) error {
	query := `
        INSERT INTO group_members (group_id, user_id, joined_at) 
        VALUES ($1, $2, CURRENT_TIMESTAMP) 
        ON CONFLICT (group_id, user_id) DO NOTHING`

	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	return err
}

func (r *GroupRepository) RemoveMember(ctx context.Context, groupID, userID uuid.UUID) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`

	_, err := r.db.ExecContext(ctx, query, groupID, userID)
	return err
}

func (r *GroupRepository) IsMember(ctx context.Context, groupID, userID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = $1 AND user_id = $2)`

	err := r.db.GetContext(ctx, &exists, query, groupID, userID)
	return exists, err
}

func (r *GroupRepository) GetByID(ctx context.Context, groupID uuid.UUID) (*models.Group, error) {
	var group models.Group
	query := `SELECT id, name, created_at FROM groups WHERE id = $1`

	err := r.db.GetContext(ctx, &group, query, groupID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrGroupNotFound.Wrap(err)
		}
		return nil, domain.NewInternalError("error fetching group by id", err)
	}

	return &group, nil
}

func (r *GroupRepository) GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]models.User, error) {
	var members []models.User
	query := `
        SELECT u.id, u.name, u.surname, u.email 
        FROM users u
        JOIN group_members gm ON u.id = gm.user_id
        WHERE gm.group_id = $1`

	err := r.db.SelectContext(ctx, &members, query, groupID)
	return members, err
}

func (r *GroupRepository) GetGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]models.Group, error) {
	var groups []models.Group
	query := `
        SELECT g.id, g.name, g.created_at 
        FROM groups g
        JOIN group_members gm ON g.id = gm.group_id
        WHERE gm.user_id = $1
        ORDER BY g.created_at DESC`

	err := r.db.SelectContext(ctx, &groups, query, userID)
	return groups, err
}
