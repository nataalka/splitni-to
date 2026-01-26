package services

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (s *UserService) Register(ctx context.Context, req models.CreateUserRequest) (*models.User, error) {
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

func (s *UserService) Login(ctx context.Context, req models.LoginRequest, jwtSecret string) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.NewInvalidCredentialsError("invalid email or password")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, domain.NewInvalidCredentialsError("wrong password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, domain.NewInternalError("failed to generate token", err)
	}

	return &models.LoginResponse{
		User:        user,
		AccessToken: tokenString,
	}, nil
}

func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
