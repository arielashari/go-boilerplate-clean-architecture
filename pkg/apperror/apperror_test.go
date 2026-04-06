package apperror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/apperror"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := apperror.New(apperror.CodeNotFound, "resource not found")

	assert.Equal(t, apperror.CodeNotFound, err.Code)
	assert.Equal(t, "resource not found", err.Message)
	assert.Nil(t, err.Internal)
	assert.Empty(t, err.Operation)
}

func TestAppError_Error(t *testing.T) {
	t.Run("without internal", func(t *testing.T) {
		err := apperror.New(apperror.CodeNotFound, "resource not found")
		assert.Equal(t, "[NOT_FOUND] resource not found", err.Error())
	})

	t.Run("with internal", func(t *testing.T) {
		internal := errors.New("db connection failed")
		err := apperror.New(apperror.CodeInternal, "internal error").WithInternal(internal)
		assert.Contains(t, err.Error(), "[INTERNAL_ERROR]")
		assert.Contains(t, err.Error(), "db connection failed")
	})
}

func TestAppError_WithInternal(t *testing.T) {
	internal := errors.New("original error")
	err := apperror.New(apperror.CodeInternal, "something failed").WithInternal(internal)

	assert.Equal(t, internal, err.Internal)
	// original must not be mutated
	assert.NotSame(t, apperror.New(apperror.CodeInternal, "something failed"), err)
}

func TestAppError_WithOperation(t *testing.T) {
	err := apperror.New(apperror.CodeNotFound, "not found").WithOperation("UserRepo.GetByID")

	assert.Equal(t, "UserRepo.GetByID", err.Operation)
	// original must not be mutated
	assert.Empty(t, apperror.New(apperror.CodeNotFound, "not found").Operation)
}

func TestAppError_Is(t *testing.T) {
	t.Run("same code matches", func(t *testing.T) {
		err1 := apperror.New(apperror.CodeNotFound, "not found")
		err2 := apperror.New(apperror.CodeNotFound, "different message")
		assert.ErrorIs(t, err1, err2)
	})

	t.Run("different code does not match", func(t *testing.T) {
		err1 := apperror.New(apperror.CodeNotFound, "not found")
		err2 := apperror.New(apperror.CodeInternal, "internal")
		assert.NotErrorIs(t, err1, err2)
	})

	t.Run("non AppError does not match", func(t *testing.T) {
		err1 := apperror.New(apperror.CodeNotFound, "not found")
		err2 := errors.New("plain error")
		assert.NotErrorIs(t, err1, err2)
	})
}

func TestAppError_Unwrap(t *testing.T) {
	internal := errors.New("wrapped")
	err := apperror.New(apperror.CodeInternal, "outer").WithInternal(internal)

	assert.Equal(t, internal, errors.Unwrap(err))
}

func TestAs(t *testing.T) {
	t.Run("succeeds for AppError", func(t *testing.T) {
		original := apperror.New(apperror.CodeForbidden, "forbidden")
		wrapped := fmt.Errorf("wrapped: %w", original)

		appErr, ok := apperror.As(wrapped)
		assert.True(t, ok)
		assert.Equal(t, apperror.CodeForbidden, appErr.Code)
	})

	t.Run("fails for non AppError", func(t *testing.T) {
		plain := errors.New("plain error")

		appErr, ok := apperror.As(plain)
		assert.False(t, ok)
		assert.Nil(t, appErr)
	})
}
