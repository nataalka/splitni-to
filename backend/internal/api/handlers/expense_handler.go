package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nataalka/splitni-to/backend/internal/domain"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
	"github.com/nataalka/splitni-to/backend/internal/domain/services"
)

type ExpenseHandler struct {
	service *services.ExpenseService
}

func NewExpenseHandler(s *services.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{service: s}
}

// Create godoc
// @Summary      Create a new expense
// @Description  Adds a new expense to a group. Sum of splits must equal the total amount.
// @Tags         expenses
// @Accept       json
// @Produce      json
// @Param        group_id path string true "Group ID"
// @Param        expense body models.CreateExpenseRequest true "Expense object"
// @Success      201 {object} models.Expense
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /groups/{group_id}/expenses [post]
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseUUID(r, "group_id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	var req models.CreateExpenseRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	expense := req.ToDomain(groupID)

	err = h.service.CreateExpense(r.Context(), expense)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, expense)
}

// ListByGroup godoc
// @Summary      List group expenses
// @Description  Returns a list of all expenses in a group (with individual splits)
// @Tags         expenses
// @Produce      json
// @Param        group_id path string true "Group ID"
// @Success      200 {array} models.Expense
// @Failure      400 {object} ErrorResponse
// @Router       /groups/{group_id}/expenses [get]
func (h *ExpenseHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseUUID(r, "group_id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	expenses, err := h.service.GetGroupExpenses(r.Context(), groupID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, expenses)
}

// GetGroupBalances godoc
// @Summary      Get group balances
// @Description  Calculates the current balance for each user in the group (who owes whom)
// @Tags         expenses
// @Produce      json
// @Param        group_id path string true "Group ID"
// @Success      200 {object} models.MemberBalance[]
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /groups/{group_id}/balances [get]
func (h *ExpenseHandler) GetGroupBalances(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseUUID(r, "group_id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	balances, err := h.service.GetGroupBalances(r.Context(), groupID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, balances)
}

// GetTotal godoc
// @Summary      Get total group spent
// @Description  Calculates the total spending in the group
// @Tags         expenses
// @Produce      json
// @Param        group_id path string true "Group ID"
// @Success      200 {object} models.TotalSpent
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /groups/{group_id}/expenses/total [get]
func (h *ExpenseHandler) GetTotal(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseUUID(r, "group_id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	total, err := h.service.GetTotalSpent(r.Context(), groupID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, models.TotalSpent{TotalSpent: total})
}

// GetUserBalance godoc
// @Summary      Get specific user balance in group
// @Description  Calculates the total expenses and awaiting incomes of the user in the group
// @Tags         expenses
// @Produce      json
// @Param        group_id path string true "Group ID"
// @Param        user_id path string true "User ID"
// @Success      200 {object} models.MemberBalance
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /groups/{group_id}/balances/{user_id} [get]
func (h *ExpenseHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseUUID(r, "group_id")
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid group id", err))
	}
	userID, err := parseUUID(r, "user_id")
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid user id"))
		return
	}

	balance, err := h.service.GetUserBalance(r.Context(), groupID, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, balance)
}

// GetExpenseDetailed godoc
// @Summary      Get expense details
// @Description  Returns detailed information about an expense including payer, group and split details with user objects.
// @Tags         expenses
// @Produce      json
// @Param        id path string true "Expense ID"
// @Success      200 {object} models.ExpenseDetailed
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /expenses/{id} [get]
func (h *ExpenseHandler) GetExpenseDetailed(w http.ResponseWriter, r *http.Request) {
	expenseID, err := parseUUID(r, "expense_id")
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid expense id", err))
		return
	}

	expense, err := h.service.GetExpenseDetailed(r.Context(), expenseID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, expense)
}
