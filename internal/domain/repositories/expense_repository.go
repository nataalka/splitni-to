package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/internal/domain/models"
	"github.com/shopspring/decimal"
)

type ExpenseRepository interface {
	CreateWithSplits(ctx context.Context, e *models.Expense) error
	GetByGroup(ctx context.Context, groupID uuid.UUID) ([]models.Expense, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Expense, error)
	GetBalancesData(ctx context.Context, groupID uuid.UUID) (map[uuid.UUID]decimal.Decimal, error)
}
