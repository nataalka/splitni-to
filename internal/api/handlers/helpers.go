package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/nataalka/splitni-to/internal/domain"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			log.Printf("failed to encode json: %v", err)
		}
	}
}

func RespondWithError(w http.ResponseWriter, err error) {
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
		case domain.TypeAuth:
			status = http.StatusUnauthorized
		case domain.TypeInternal:
			status = http.StatusInternalServerError
		}
	} else {
		log.Printf("[UNHANDLED ERROR]: %v", err)
	}

	WriteJSON(w, status, resp)
}
