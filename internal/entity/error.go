package entity

import "github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/apperror"

var (
	ErrNotFound           = apperror.New(apperror.CodeNotFound, "resource not found")
	ErrConflict           = apperror.New(apperror.CodeConflict, "resource already exists")
	ErrInvalidInput       = apperror.New(apperror.CodeValidation, "invalid input provided")
	ErrUnauthorized       = apperror.New(apperror.CodeUnauthorized, "unauthorized access")
	ErrForbidden          = apperror.New(apperror.CodeForbidden, "permission denied")
	ErrInvalidCredentials = apperror.New(apperror.CodeInvalidCreds, "invalid email or password")
	ErrEmailAlreadyExists = apperror.New(apperror.CodeEmailTaken, "email already exists")
	ErrInternal           = apperror.New(apperror.CodeInternal, "internal server error")
)
