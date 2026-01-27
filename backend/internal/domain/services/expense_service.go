package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/backend/internal/domain"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
	"github.com/nataalka/splitni-to/backend/internal/domain/repositories"
	"github.com/shopspring/decimal"
)

type ExpenseService struct {
	repo      repositories.ExpenseRepository
	userRepo  repositories.UserRepository
	groupRepo repositories.GroupRepository
}

func NewExpenseService(e repositories.ExpenseRepository, u repositories.UserRepository, g repositories.GroupRepository) *ExpenseService {
	return &ExpenseService{
		repo:      e,
		userRepo:  u,
		groupRepo: g,
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

func (s *ExpenseService) GetExpenseDetailed(ctx context.Context, expenseID uuid.UUID) (*models.ExpenseDetailed, error) {
	expense, err := s.repo.GetByID(ctx, expenseID)
	if err != nil {
		return nil, domain.ErrExpenseNotFound
	}

	group, err := s.groupRepo.GetByID(ctx, expense.GroupID)
	if err != nil {
		return nil, domain.ErrGroupNotFound
	}

	userIDs := []uuid.UUID{expense.PayerID}
	for _, split := range expense.Splits {
		userIDs = append(userIDs, split.UserID)
	}

	users, err := s.userRepo.GetByIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	userMap := make(map[uuid.UUID]models.User)
	for _, u := range users {
		userMap[u.ID] = u
	}

	detailedSplits := make([]models.ExpenseSplitDetailed, 0, len(expense.Splits))
	for _, split := range expense.Splits {
		detailedSplits = append(detailedSplits, models.ExpenseSplitDetailed{
			User:   userMap[split.UserID],
			Amount: split.Amount,
		})
	}

	return &models.ExpenseDetailed{
		ID:          expense.ID,
		Group:       *group,
		Payer:       userMap[expense.PayerID],
		Amount:      expense.Amount,
		Currency:    expense.Currency,
		Description: expense.Description,
		CreatedAt:   expense.CreatedAt,
		Splits:      detailedSplits,
	}, nil
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

func (s *ExpenseService) GetTotalSpent(ctx context.Context, groupID uuid.UUID) (decimal.Decimal, error) {
	total, err := s.repo.GetTotalByGroup(ctx, groupID)
	if err != nil {
		return decimal.Zero, domain.NewInternalError("failed to calculate total spent", err)
	}
	return total, nil
}

func (s *ExpenseService) GetUserBalance(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (*models.MemberBalance, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	balance, err := s.repo.GetUserBalance(ctx, groupID, userID)
	if err != nil {
		return nil, domain.NewInternalError("failed to calculate user balance", err)
	}

	return &models.MemberBalance{
		User:    *user,
		Balance: balance,
	}, nil
}
