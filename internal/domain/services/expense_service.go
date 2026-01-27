package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain"
	"github.com/nataalka/splitni-to/internal/domain/models"
	"github.com/nataalka/splitni-to/internal/domain/repositories"
	"github.com/shopspring/decimal"
)

type ExpenseService struct {
	repo     repositories.ExpenseRepository
	userRepo repositories.UserRepository
}

func NewExpenseService(e repositories.ExpenseRepository, u repositories.UserRepository) *ExpenseService {
	return &ExpenseService{
		repo:     e,
		userRepo: u,
	}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, e *models.Expense) error {
	if e.Amount.LessThanOrEqual(decimal.Zero) {
		return domain.ErrInvalidAmount
	}

	if len(e.Splits) == 0 {
		return domain.ErrNoSplits
	}

	totalSplits := decimal.Zero
	for _, split := range e.Splits {
		totalSplits = totalSplits.Add(split.Amount)
	}

	if !totalSplits.Equal(e.Amount) {
		return domain.ErrInvalidSplits
	}

	e.ID = uuid.New()
	e.CreatedAt = time.Now()

	for i := range e.Splits {
		e.Splits[i].ExpenseID = e.ID
	}

	return s.repo.CreateWithSplits(ctx, e)
}

func (s *ExpenseService) GetGroupExpenses(ctx context.Context, groupID uuid.UUID) ([]models.Expense, error) {
	expenses, err := s.repo.GetByGroup(ctx, groupID)
	if err != nil {
		return nil, domain.NewInternalError("failed to fetch expenses", err)
	}
	return expenses, nil
}

func (s *ExpenseService) GetExpenseDetails(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	expense, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrExpenseNotFound
	}
	return expense, nil
}

func (s *ExpenseService) GetGroupBalances(ctx context.Context, groupID uuid.UUID) ([]models.MemberBalance, error) {

	rawBalances, err := s.repo.GetBalancesData(ctx, groupID)
	if err != nil {
		return nil, domain.NewInternalError("failed to get balance data", err)
	}

	var userIDs []uuid.UUID
	for id := range rawBalances {
		userIDs = append(userIDs, id)
	}

	users, err := s.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, domain.NewInternalError("failed to fetch user details", err)
	}

	var result []models.MemberBalance

	userMap := make(map[uuid.UUID]models.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	for uID, bal := range rawBalances {
		user, ok := userMap[uID]
		if !ok {
			continue
		}

		result = append(result, models.MemberBalance{
			User:    user,
			Balance: bal,
		})
	}

	return result, nil
}
