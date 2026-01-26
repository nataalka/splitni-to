package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain"
	"github.com/nataalka/splitni-to/internal/domain/models"
	"github.com/nataalka/splitni-to/internal/domain/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, req models.CreateUser) (*models.User, error) {
	if req.Email == "" || req.Password == "" {
		return nil, domain.NewValidationError("email and password are required")
	}

	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, domain.ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, domain.NewInternalError("failed to hash password", err)
	}

	user := &models.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Surname:      req.Surname,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now().UTC(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
