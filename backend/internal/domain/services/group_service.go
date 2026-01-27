package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/backend/internal/domain"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
	"github.com/nataalka/splitni-to/backend/internal/infrastructure/postgres"
)

type GroupService struct {
	repo *postgres.GroupRepository
}

func NewGroupService(repo *postgres.GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func (s *GroupService) CreateGroup(ctx context.Context, req models.CreateGroupRequest, creatorID uuid.UUID) (*models.Group, error) {
	group := &models.Group{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
	}

	err := s.repo.Create(ctx, group, creatorID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *GroupService) UpdateGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID, req models.CreateGroupRequest) error {
	isMember, err := s.repo.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return domain.ErrNotInGroup
	}

	group := &models.Group{
		ID:          groupID,
		Name:        req.Name,
		Description: req.Description,
	}

	return s.repo.Update(ctx, group)
}

func (s *GroupService) DeleteGroup(ctx context.Context, userID uuid.UUID, groupID uuid.UUID) error {
	isMember, err := s.repo.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return domain.ErrNotInGroup
	}

	return s.repo.Delete(ctx, groupID)
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
