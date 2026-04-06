package apperror

import (
	"errors"
	"fmt"
)

type ErrorCode string

const (
	CodeValidation   ErrorCode = "VALIDATION_ERROR"
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	CodeForbidden    ErrorCode = "FORBIDDEN"
	CodeInvalidCreds ErrorCode = "INVALID_CREDENTIALS"
	CodeNotFound     ErrorCode = "NOT_FOUND"
	CodeConflict     ErrorCode = "CONFLICT"
	CodeEmailTaken   ErrorCode = "EMAIL_TAKEN"
	CodeInternal     ErrorCode = "INTERNAL_ERROR"
	CodeUnavailable  ErrorCode = "SERVICE_UNAVAILABLE"
	CodeTimeout      ErrorCode = "TIMEOUT"
)

type AppError struct {
	Code      ErrorCode `json:"code"`
	Message   string    `json:"message"`
	Internal  error     `json:"-"`
	Operation string    `json:"-"`
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Internal
}

func (e *AppError) Is(target error) bool {
	other, ok := target.(*AppError)
	if !ok || other == nil {
		return false
	}
	return e.Code == other.Code
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func (e *AppError) WithInternal(err error) *AppError {
	newErr := *e
	newErr.Internal = err
	return &newErr
}

func (e *AppError) WithOperation(op string) *AppError {
	newErr := *e
	newErr.Operation = op
	return &newErr
}

func As(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}
