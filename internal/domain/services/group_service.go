package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain/models"
	"github.com/nataalka/splitni-to/internal/infrastructure/postgres"
)

type GroupService struct {
	repo *postgres.GroupRepository
}

func NewGroupService(repo *postgres.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) CreateGroup(ctx context.Context, name string, creatorID uuid.UUID) (*models.Group, error) {
	group := &models.Group{
		ID:        uuid.New(),
		Name:      name,
		CreatedAt: time.Now(),
	}

	err := s.repo.Create(ctx, group, creatorID)
	if err != nil {
		return nil, err
	}

	return group, nil
}
