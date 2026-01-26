package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nataalka/splitni-to/internal/domain/models"
)

type FriendRepository struct {
	db *sqlx.DB
}

func NewPostgresFriendRepository(db *sqlx.DB) *FriendRepository {
	return &FriendRepository{db: db}
}

func (r *FriendRepository) CreateFriendRequest(ctx context.Context, id1, id2 uuid.UUID) error {
	uid1, uid2 := sortUUIDs(id1, id2)
	_, err := r.db.ExecContext(ctx, "INSERT INTO friendships (user_id1, user_id2, status) VALUES ($1, $2, 'PENDING') ON CONFLICT DO NOTHING", uid1, uid2)
	return err
}

func (r *FriendRepository) AcceptFriendRequest(ctx context.Context, id1, id2 uuid.UUID) error {
	uid1, uid2 := sortUUIDs(id1, id2)
	_, err := r.db.ExecContext(ctx, "UPDATE friendships SET status = 'ACCEPTED' WHERE user_id1 = $1 AND user_id2 = $2 AND status = 'PENDING'", uid1, uid2)
	return err
}

func (r *FriendRepository) DeleteFriendship(ctx context.Context, id1, id2 uuid.UUID) error {
	uid1, uid2 := sortUUIDs(id1, id2)
	_, err := r.db.ExecContext(ctx, "DELETE FROM friendships WHERE user_id1 = $1 AND user_id2 = $2", uid1, uid2)
	return err
}

func (r *FriendRepository) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]models.User, error) {
	var pending []models.User
	query := `
        SELECT u.id, u.name, u.surname, u.email 
        FROM users u
        JOIN friendships f ON (f.user_id1 = u.id OR f.user_id2 = u.id)
        WHERE (f.user_id1 = $1 OR f.user_id2 = $1) 
          AND u.id != $1 
          AND f.status = 'PENDING'`

	err := r.db.SelectContext(ctx, &pending, query, userID)
	return pending, err
}

func (r *FriendRepository) GetFriends(ctx context.Context, userID uuid.UUID) ([]models.User, error) {
	var friends []models.User
	query := `
        SELECT u.id, u.name, u.surname, u.email 
        FROM users u
        JOIN friendships f ON (f.user_id1 = u.id OR f.user_id2 = u.id)
        WHERE (f.user_id1 = $1 OR f.user_id2 = $1) 
          AND u.id != $1 
          AND f.status = 'ACCEPTED'`

	err := r.db.SelectContext(ctx, &friends, query, userID)
	return friends, err
}

func sortUUIDs(id1, id2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if id1.String() < id2.String() {
		return id1, id2
	}
	return id2, id1
}
