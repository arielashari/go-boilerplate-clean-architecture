package validator_test

import (
	"testing"

	customervalidator "github.com/Primuse-Pte-Ltd/go-boilerplate-clean-architecture/pkg/validator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=3"`
	Phone string `json:"phone" validate:"required,numeric"`
}

func TestGetValidator(t *testing.T) {
	v := customervalidator.GetValidator()
	assert.NotNil(t, v)
}

func TestGetValidator_UsesJSONFieldNames(t *testing.T) {
	v := customervalidator.GetValidator()

	type s struct {
		EmailAddress string `json:"email_address" validate:"required"`
	}

	err := v.Struct(s{})
	assert.Error(t, err)

	var valErrs validator.ValidationErrors
	assert.ErrorAs(t, err, &valErrs)
	// should use json tag name not struct field name
	assert.Equal(t, "email_address", valErrs[0].Field())
}

func TestFormatValidationError(t *testing.T) {
	v := customervalidator.GetValidator()

	t.Run("required", func(t *testing.T) {
		type s struct {
			Name string `json:"name" validate:"required"`
		}
		err := v.Struct(s{})
		var valErrs validator.ValidationErrors
		assert.ErrorAs(t, err, &valErrs)

		result := customervalidator.FormatValidationError(valErrs)
		assert.Len(t, result, 1)
		assert.Equal(t, "name", result[0].Field)
		assert.Equal(t, "is required", result[0].Message)
	})

	t.Run("email", func(t *testing.T) {
		type s struct {
			Email string `json:"email" validate:"required,email"`
		}
		err := v.Struct(s{Email: "not-an-email"})
		var valErrs validator.ValidationErrors
		assert.ErrorAs(t, err, &valErrs)

		result := customervalidator.FormatValidationError(valErrs)
		assert.Len(t, result, 1)
		assert.Equal(t, "email", result[0].Field)
		assert.Equal(t, "must be a valid email address", result[0].Message)
	})

	t.Run("min", func(t *testing.T) {
		type s struct {
			Name string `json:"name" validate:"required,min=3"`
		}
		err := v.Struct(s{Name: "ab"})
		var valErrs validator.ValidationErrors
		assert.ErrorAs(t, err, &valErrs)

		result := customervalidator.FormatValidationError(valErrs)
		assert.Len(t, result, 1)
		assert.Equal(t, "name", result[0].Field)
		assert.Equal(t, "must be at least 3 characters", result[0].Message)
	})

	t.Run("numeric", func(t *testing.T) {
		type s struct {
			Phone string `json:"phone" validate:"required,numeric"`
		}
		err := v.Struct(s{Phone: "not-numeric"})
		var valErrs validator.ValidationErrors
		assert.ErrorAs(t, err, &valErrs)

		result := customervalidator.FormatValidationError(valErrs)
		assert.Len(t, result, 1)
		assert.Equal(t, "phone", result[0].Field)
		assert.Equal(t, "must be a numeric value", result[0].Message)
	})

	t.Run("unknown tag falls back to default message", func(t *testing.T) {
		type s struct {
			Age int `json:"age" validate:"min=18"`
		}
		err := v.Struct(s{Age: 10})
		var valErrs validator.ValidationErrors
		assert.ErrorAs(t, err, &valErrs)

		result := customervalidator.FormatValidationError(valErrs)
		assert.Len(t, result, 1)
		assert.Equal(t, "age", result[0].Field)
	})

	t.Run("multiple errors", func(t *testing.T) {
		err := v.Struct(testStruct{})
		var valErrs validator.ValidationErrors
		assert.ErrorAs(t, err, &valErrs)

		result := customervalidator.FormatValidationError(valErrs)
		assert.Len(t, result, 3)
	})

	t.Run("no errors returns nil", func(t *testing.T) {
		result := customervalidator.FormatValidationError(nil)
		assert.Nil(t, result)
	})
}
