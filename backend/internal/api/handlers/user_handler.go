package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nataalka/splitni-to/backend/internal/domain"
	"github.com/nataalka/splitni-to/backend/internal/domain/models"
	"github.com/nataalka/splitni-to/backend/internal/domain/services"
)

type UserHandler struct {
	userService *services.UserService
	jwtSecret   string
}

func NewUserHandler(us *services.UserService, jwtSecret string) *UserHandler {
	return &UserHandler{
		userService: us,
		jwtSecret:   jwtSecret,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account and return user data
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body models.CreateUserRequest true "Registration details"
// @Success      201 {object} models.User
// @Failure      400 {object} domain.AppError
// @Router       /register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	user, err := h.userService.Register(ctx, req)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body models.LoginRequest true "Login credentials"
// @Success      200 {object} models.LoginResponse
// @Failure      401 {object} domain.AppError
// @Router       /login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, domain.NewValidationError("invalid request body", err))
		return
	}

	if req.Email == "" || req.Password == "" {
		respondWithError(w, domain.NewValidationError("email and password are required"))
		return
	}

	res, err := h.userService.Login(ctx, req, h.jwtSecret)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, res)
}

// GetMe godoc
// @Summary      Get current user profile
// @Description  Fetch information of the currently authenticated user
// @Tags         auth
// @Security     BearerAuth
// @Produce      json
// @Success      200 {object} models.User
// @Failure      401 {object} domain.AppError
// @Router       /me [get]
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		respondWithError(w, domain.NewAuthError("unauthorized: user id not found in context"))
		return
	}

	user, err := h.userService.GetByID(ctx, userID)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// GetByID godoc
// @Summary      Get user profile
// @Description  Fetch user information by their unique ID
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        user_id path string true "User UUID"
// @Success      200 {object} models.User
// @Failure      404 {object} domain.AppError
// @Router       /users/{user_id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := parseUUID(r, "user_id")
	if err != nil {
		respondWithError(w, err)
		return
	}

	user, err := h.userService.GetByID(ctx, id)
	if err != nil {
		respondWithError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
