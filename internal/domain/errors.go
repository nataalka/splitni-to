package domain

import (
	"fmt"
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Wrap(err error) *AppError {
	e.Err = err
	return e
}

type ErrorType string

const (
	TypeValidation         ErrorType = "VALIDATION_ERROR"
	TypeNotFound           ErrorType = "NOT_FOUND"
	TypeConflict           ErrorType = "CONFLICT"
	TypeInternal           ErrorType = "INTERNAL_ERROR"
	TypeAuth               ErrorType = "UNAUTHORIZED"
	TypeInvalidCredentials ErrorType = "INVALID_CREDENTIALS"
	TypePermission         ErrorType = "PERMISSION_DENIED"
)

func NewValidationError(msg string, errs ...error) *AppError {
	var err error
	if len(errs) > 0 {
		err = errs[0]
	}
	return &AppError{Type: TypeValidation, Message: msg, Err: err}
}

func NewNotFoundError(msg string) *AppError {
	return &AppError{Type: TypeNotFound, Message: msg}
}

func NewConflictError(msg string) *AppError {
	return &AppError{Type: TypeConflict, Message: msg}
}

func NewInternalError(msg string, err error) *AppError {
	return &AppError{Type: TypeInternal, Message: msg, Err: err}
}

func NewAuthError(msg string) *AppError {
	return &AppError{Type: TypeAuth, Message: msg}
}

func NewInvalidCredentialsError(msg string) *AppError {
	return &AppError{Type: TypeInvalidCredentials, Message: msg}
}

func NewPermissionError(msg string) *AppError {
	return &AppError{Type: TypePermission, Message: msg}
}

var (
	ErrEmailTaken     = NewConflictError("email is already registered")
	ErrUserNotFound   = NewNotFoundError("user not found")
	ErrGroupNotFound  = NewNotFoundError("group not found")
	ErrAlreadyInGroup = NewConflictError("member is already in group")
	ErrNotInGroup     = NewNotFoundError("member is not in group")
)
