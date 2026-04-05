package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "" || name == "-" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}
		return name
	})
}
func GetValidator() *validator.Validate {
	return validate
}

func FormatValidationError(err error) []FieldError {
	var errors []FieldError
	if castedObject, ok := err.(validator.ValidationErrors); ok {
		for _, e := range castedObject {
			field := e.Field()
			tag := e.Tag()
			param := e.Param()

			var message string
			switch tag {
			case "required":
				message = "is required"
			case "email":
				message = "must be a valid email address"
			case "min":
				message = fmt.Sprintf("must be at least %s characters", param)
			case "numeric":
				message = "must be a numeric value"
			default:
				message = fmt.Sprintf("failed on the '%s' validation", tag)
			}

			errors = append(errors, FieldError{
				Field:   field,
				Message: message,
			})
		}
	}
	return errors
}
