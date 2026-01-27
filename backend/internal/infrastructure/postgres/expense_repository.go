package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
	"github.com/shopspring/decimal"
)

type ExpenseRepository struct {
	db *sqlx.DB
}

func NewExpenseRepository(db *sqlx.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) CreateWithSplits(ctx context.Context, e *models.Expense) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO expenses (id, group_id, payer_id, amount, currency, description, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.ExecContext(ctx, query, e.ID, e.GroupID, e.PayerID, e.Amount, e.Currency, e.Description, e.CreatedAt)
	if err != nil {
		return err
	}

	splitQuery := `INSERT INTO expense_splits (expense_id, user_id, amount) VALUES ($1, $2, $3)`
	for _, split := range e.Splits {
		_, err = tx.ExecContext(ctx, splitQuery, e.ID, split.UserID, split.Amount)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *ExpenseRepository) GetByGroup(ctx context.Context, groupID uuid.UUID) ([]models.Expense, error) {
	query := `SELECT id, group_id, payer_id, amount, currency, description, created_at 
              FROM expenses WHERE group_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		err = rows.Scan(&e.ID, &e.GroupID, &e.PayerID, &e.Amount, &e.Currency, &e.Description, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, e)
	}
	return expenses, nil
}

func (r *ExpenseRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	var e models.Expense

	query := `SELECT id, group_id, payer_id, amount, currency, description, created_at 
              FROM expenses WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&e.ID, &e.GroupID, &e.PayerID, &e.Amount, &e.Currency, &e.Description, &e.CreatedAt)
	if err != nil {
		return nil, err
	}

	splitQuery := `SELECT expense_id, user_id, amount FROM expense_splits WHERE expense_id = $1`
	rows, err := r.db.QueryContext(ctx, splitQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.ExpenseSplit
		err = rows.Scan(&s.ExpenseID, &s.UserID, &s.Amount)
		if err != nil {
			return nil, err
		}
		e.Splits = append(e.Splits, s)
	}

	return &e, nil
}

func (r *ExpenseRepository) GetBalancesData(ctx context.Context, groupID uuid.UUID) (map[uuid.UUID]decimal.Decimal, error) {
	query := `
       SELECT user_id, SUM(diff) as balance
       FROM (
          SELECT payer_id as user_id, amount as diff 
          FROM expenses WHERE group_id = $1
          UNION ALL
          SELECT s.user_id, -s.amount as diff 
          FROM expense_splits s
          JOIN expenses e ON s.expense_id = e.id
          WHERE e.group_id = $1
       ) combined
       GROUP BY user_id`

	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make(map[uuid.UUID]decimal.Decimal)
	for rows.Next() {
		var uID uuid.UUID
		var bal decimal.Decimal
		err = rows.Scan(&uID, &bal)
		if err != nil {
			return nil, err
		}
		balances[uID] = bal
	}
	return balances, nil
}
