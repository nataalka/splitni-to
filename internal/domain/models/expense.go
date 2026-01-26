package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Expense struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	GroupID     uuid.UUID       `json:"group_id" db:"group_id"`
	PayerID     uuid.UUID       `json:"payer_id" db:"payer_id"`
	Amount      decimal.Decimal `json:"amount" db:"amount"`
	Description string          `json:"description" db:"description"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	Splits      []ExpenseSplit  `json:"splits,omitempty"`
}

type ExpenseSplit struct {
	ExpenseID uuid.UUID       `json:"expense_id" postgres:"expense_id"`
	UserID    uuid.UUID       `json:"user_id" postgres:"user_id"`
	Amount    decimal.Decimal `json:"amount" postgres:"amount"`
}
