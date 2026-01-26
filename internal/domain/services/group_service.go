package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain"
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

func (s *GroupService) AddMemberToGroup(ctx context.Context, actorID, groupID, newMemberID uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, groupID)
	if err != nil {
		return err
	}

	isMember, err := s.repo.IsMember(ctx, groupID, actorID)
	if err != nil {
		return err
	}
	if !isMember {
		return domain.ErrNotInGroup
	}

	isAlreadyThere, _ := s.repo.IsMember(ctx, groupID, newMemberID)
	if isAlreadyThere {
		return domain.ErrAlreadyInGroup
	}

	return s.repo.AddMember(ctx, groupID, newMemberID)
}

func (s *GroupService) RemoveMember(ctx context.Context, actorID, groupID, targetID uuid.UUID) error {
	isMember, err := s.repo.IsMember(ctx, groupID, actorID)
	if err != nil {
		return err
	}
	if !isMember {
		return domain.ErrNotInGroup
	}

	return s.repo.RemoveMember(ctx, groupID, targetID)
}

func (s *GroupService) GetByID(ctx context.Context, userID, groupID uuid.UUID) (*models.Group, error) {
	isMember, err := s.repo.IsMember(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, domain.ErrNotInGroup
	}

	return s.repo.GetByID(ctx, groupID)
}

func (s *GroupService) GetMembers(ctx context.Context, userID, groupID uuid.UUID) ([]models.User, error) {
	isMember, err := s.repo.IsMember(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, domain.ErrNotInGroup
	}

	return s.repo.GetGroupMembers(ctx, groupID)
}

func (s *GroupService) GetUserGroups(ctx context.Context, userID uuid.UUID) ([]models.Group, error) {
	return s.repo.GetGroupsByUserID(ctx, userID)
}
