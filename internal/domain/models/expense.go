package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Expense struct {
	ID          uuid.UUID       `json:"id" postgres:"id"`
	GroupID     uuid.UUID       `json:"group_id" postgres:"group_id"`
	PayerID     uuid.UUID       `json:"payer_id" postgres:"payer_id"`
	Amount      decimal.Decimal `json:"amount" postgres:"amount"`
	Description string          `json:"description" postgres:"description"`
	CreatedAt   time.Time       `json:"created_at" postgres:"created_at"`
	Splits      []ExpenseSplit  `json:"splits,omitempty"`
}

type ExpenseSplit struct {
	ExpenseID uuid.UUID       `json:"expense_id" postgres:"expense_id"`
	UserID    uuid.UUID       `json:"user_id" postgres:"user_id"`
	Amount    decimal.Decimal `json:"amount" postgres:"amount"`
}
