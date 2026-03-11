package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nataalka/splitni-to/backend/internal/domain"
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

func (r *ExpenseRepository) Update(ctx context.Context, expense *models.Expense, splits []models.ExpenseSplit) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return domain.NewInternalError("failed to begin transaction", err)
	}
	defer tx.Rollback()

	query := `UPDATE expenses SET description = $1, amount = $2, payer_id = $3 WHERE id = $4`
	_, err = tx.ExecContext(ctx, query, expense.Description, expense.Amount, expense.PayerID, expense.ID)
	if err != nil {
		return domain.NewInternalError("could not update expense", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM expense_splits WHERE expense_id = $1", expense.ID)
	if err != nil {
		return domain.NewInternalError("could not clear old splits", err)
	}

	for _, s := range splits {
		_, err = tx.ExecContext(ctx,
			"INSERT INTO expense_splits (expense_id, user_id, amount) VALUES ($1, $2, $3)",
			expense.ID, s.UserID, s.Amount)
		if err != nil {
			return domain.NewInternalError("could not insert new splits", err)
		}
	}

	return tx.Commit()
}

func (r *ExpenseRepository) Delete(ctx context.Context, expenseID uuid.UUID) error {
	query := `DELETE FROM expenses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, expenseID)
	if err != nil {
		return domain.NewInternalError("could not delete expense", err)
	}
	return nil
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
	var expenseIDs []uuid.UUID

	for rows.Next() {
		var e models.Expense
		err = rows.Scan(&e.ID, &e.GroupID, &e.PayerID, &e.Amount, &e.Currency, &e.Description, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		e.Splits = []models.ExpenseSplit{}
		expenses = append(expenses, e)
		expenseIDs = append(expenseIDs, e.ID)
	}

	if len(expenses) == 0 {
		return expenses, nil
	}

	querySplits := `SELECT expense_id, user_id, amount FROM expense_splits WHERE expense_id = ANY($1)`

	splitRows, err := r.db.QueryContext(ctx, querySplits, pq.Array(expenseIDs))
	if err != nil {
		return expenses, nil
	}
	defer splitRows.Close()

	splitsByExpense := make(map[uuid.UUID][]models.ExpenseSplit)
	for splitRows.Next() {
		var s models.ExpenseSplit
		err = splitRows.Scan(&s.ExpenseID, &s.UserID, &s.Amount)
		if err != nil {
			continue
		}
		splitsByExpense[s.ExpenseID] = append(splitsByExpense[s.ExpenseID], s)
	}

	for i := range expenses {
		if s, ok := splitsByExpense[expenses[i].ID]; ok {
			expenses[i].Splits = s
		}
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

func (r *ExpenseRepository) GetTotalByGroup(ctx context.Context, groupID uuid.UUID) (decimal.Decimal, error) {
	var total decimal.Decimal
	query := `SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE group_id = $1`

	err := r.db.GetContext(ctx, &total, query, groupID)
	if err != nil {
		return decimal.Zero, err
	}

	return total, nil
}

func (r *ExpenseRepository) GetUserBalance(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (decimal.Decimal, error) {
	var balance decimal.Decimal
	query := `
		SELECT 
			(SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE group_id = $1 AND payer_id = $2) -
			(SELECT COALESCE(SUM(es.amount), 0) FROM expense_splits es 
			 JOIN expenses e ON es.expense_id = e.id 
			 WHERE e.group_id = $1 AND es.user_id = $2)
	`

	err := r.db.GetContext(ctx, &balance, query, groupID, userID)
	if err != nil {
		return decimal.Zero, err
	}

	return balance, nil
}
