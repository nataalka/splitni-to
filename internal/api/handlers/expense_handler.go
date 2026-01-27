package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nataalka/splitni-to/internal/domain"
	"github.com/nataalka/splitni-to/internal/domain/models"
	"github.com/nataalka/splitni-to/internal/domain/services"
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
// @Param        groupID path string true "Group ID"
// @Param        expense body models.CreateExpenseRequest true "Expense object"
// @Success      201 {object} models.Expense
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /groups/{groupID}/expenses [post]
func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseUUID(r, "groupID")
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
// @Description  Returns a list of all expenses in a group (without individual splits)
// @Tags         expenses
// @Produce      json
// @Param        groupID path string true "Group ID"
// @Success      200 {array} models.Expense
// @Failure      400 {object} ErrorResponse
// @Router       /groups/{groupID}/expenses [get]
func (h *ExpenseHandler) ListByGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseUUID(r, "groupID")
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
// @Param        groupID path string true "Group ID"
// @Success      200 {object} models.MemberBalance
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /groups/{groupID}/balances [get]
func (h *ExpenseHandler) GetGroupBalances(w http.ResponseWriter, r *http.Request) {
	groupID, err := parseUUID(r, "groupID")
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
