package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nataalka/splitni-to/backend/internal/domain"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewPostgresUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
       INSERT INTO users (id, name, surname, email, password_hash, created_at)
       VALUES ($1, $2, $3, $4, $5, $6)
    `
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Surname, user.Email, user.PasswordHash, user.CreatedAt)

	if err != nil {
		return domain.NewInternalError("could not create user", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, surname, email, password_hash, created_at FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound.Wrap(err)
		}
		return nil, domain.NewInternalError("error fetching user by email", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT id, name, surname, email, password_hash, created_at FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound.Wrap(err)
		}
		return nil, domain.NewInternalError("error fetching user by id", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]models.User, error) {
	if len(ids) == 0 {
		return []models.User{}, nil
	}

	query := `SELECT id, email, name, surname FROM users WHERE id = ANY($1)`

	var users []models.User
	err := r.db.SelectContext(ctx, &users, query, pq.Array(ids))
	if err != nil {
		return nil, err
	}

	return users, nil
}
