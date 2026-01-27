package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/nataalka/splitni-to/backend/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func writeJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("failed to encode json: %v", err)
		}
	}
}

func respondWithError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError

	status := http.StatusInternalServerError
	resp := ErrorResponse{
		Error: "Internal server error",
		Code:  string(domain.TypeInternal),
	}

	if errors.As(err, &appErr) {
		if appErr.Err != nil {
			log.Printf("[DEBUG] %s: %s | Origin: %v", appErr.Type, appErr.Message, appErr.Err)
		}

		resp.Error = appErr.Message
		resp.Code = string(appErr.Type)

		switch appErr.Type {
		case domain.TypeValidation:
			status = http.StatusBadRequest
		case domain.TypeConflict:
			status = http.StatusConflict
		case domain.TypeNotFound:
			status = http.StatusNotFound
		case domain.TypeAuth, domain.TypeInvalidCredentials:
			status = http.StatusUnauthorized
		case domain.TypePermission:
			status = http.StatusForbidden
		case domain.TypeInternal:
			status = http.StatusInternalServerError
		}
	} else {
		log.Printf("[UNHANDLED ERROR]: %v", err)
	}

	writeJSON(w, status, resp)
}

func getUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return uuid.Nil, domain.NewInternalError("failed to get claims", err)
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, domain.NewAuthError("user_id (sub) not found in token claims")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, domain.NewAuthError("invalid user_id format in token")
	}

	return userID, nil
}

func parseUUID(r *http.Request, paramName string) (uuid.UUID, error) {
	val := chi.URLParam(r, paramName)
	if val == "" {
		return uuid.Nil, domain.NewValidationError(fmt.Sprintf("missing parameter %s", paramName), nil)
	}

	id, err := uuid.Parse(val)
	if err != nil {
		return uuid.Nil, domain.NewValidationError(fmt.Sprintf("invalid UUID format for %s", paramName), err)
	}

	return id, nil
}
