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
	Currency    string          `json:"currency" db:"currency"`
	Description string          `json:"description" db:"description"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	Splits      []ExpenseSplit  `json:"splits,omitempty"`
}

type ExpenseSplit struct {
	ExpenseID uuid.UUID       `json:"expense_id" postgres:"expense_id"`
	UserID    uuid.UUID       `json:"user_id" postgres:"user_id"`
	Amount    decimal.Decimal `json:"amount" postgres:"amount"`
}

type CreateExpenseRequest struct {
	PayerID     uuid.UUID                   `json:"payer_id"`
	Amount      decimal.Decimal             `json:"amount"`
	Currency    string                      `json:"currency"`
	Description string                      `json:"description"`
	Splits      []CreateExpenseSplitRequest `json:"splits"`
}

type CreateExpenseSplitRequest struct {
	UserID uuid.UUID       `json:"user_id"`
	Amount decimal.Decimal `json:"amount"`
}

func (req *CreateExpenseRequest) ToDomain(groupID uuid.UUID) *Expense {
	expense := &Expense{
		GroupID:     groupID,
		PayerID:     req.PayerID,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
	}

	for _, s := range req.Splits {
		expense.Splits = append(expense.Splits, ExpenseSplit{
			UserID: s.UserID,
			Amount: s.Amount,
		})
	}
	return expense
}

type MemberBalance struct {
	User    User            `json:"user"`
	Balance decimal.Decimal `json:"balance"`
}

type TotalSpent struct {
	TotalSpent decimal.Decimal `json:"total_spent"`
}
