package booking

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	return "validation error"
}

func NewValidationError(err error) *ValidationError {
	var validationErrors validator.ValidationErrors
	var fieldErrors []FieldError

	if errors.As(err, &validationErrors) {
		for _, fieldErr := range validationErrors {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   fieldErr.Field(),
				Message: getValidationMessage(fieldErr),
			})
		}
	} else if err != nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   "general",
			Message: err.Error(),
		})
	}

	return &ValidationError{
		Fields: fieldErrors,
	}
}

func getValidationMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldErr.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldErr.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fieldErr.Field(), fieldErr.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", fieldErr.Field(), fieldErr.Param())
	default:
		return fmt.Sprintf("%s is invalid", fieldErr.Field())
	}
}
