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

// Update godoc
// @Summary      Update an existing expense
// @Description  Update the description, amount, payer, or split distribution of an expense. Only members of the group can perform this action.
// @Tags         expenses
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        expense_id path string true "Expense UUID"
// @Param        expense body models.CreateExpenseRequest true "Updated expense data"
// @Success      204 "No Content - Expense successfully updated"
// @Failure      400 {object} domain.AppError "Invalid UUID or malformed JSON body"
// @Failure      401 {object} domain.AppError "Unauthorized - User not logged in"
// @Failure      403 {object} domain.AppError "Forbidden - User is not a member of the group"
// @Failure      404 {object} domain.AppError "Expense not found"
// @Failure      500 {object} domain.AppError "Internal server error"
// @Router       /expenses/{expense_id} [patch]
func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := getUserIDFromContext(ctx)
	expenseID, _ := parseUUID(r, "expense_id")

	var req models.CreateExpenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	if err := h.service.UpdateExpense(ctx, userID, expenseID, req); err != nil {
		respondWithError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete godoc
// @Summary      Delete an expense
// @Description  Permanently remove an expense and its associated splits. Only members of the group can perform this action.
// @Tags         expenses
// @Security     BearerAuth
// @Param        expense_id path string true "Expense UUID"
// @Success      204 "No Content - Expense successfully deleted"
// @Failure      400 {object} domain.AppError "Invalid UUID"
// @Failure      401 {object} domain.AppError "Unauthorized"
// @Failure      403 {object} domain.AppError "Forbidden - User not in group"
// @Failure      404 {object} domain.AppError "Expense not found"
// @Failure      500 {object} domain.AppError "Internal server error"
// @Router       /expenses/{expense_id} [delete]
func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := getUserIDFromContext(ctx)
	expenseID, _ := parseUUID(r, "expense_id")

	if err := h.service.DeleteExpense(ctx, userID, expenseID); err != nil {
		respondWithError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
