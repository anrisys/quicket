package errs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	AppError
	Fields []FieldError `json:"fields,omitempty"`
}

func (e *ValidationError) Error() string {
	baseMsg := e.AppError.Error()
	if len(e.Fields) > 0 {
		return fmt.Sprintf("%s - Validation Details: %v", baseMsg, e.Fields)
	}
	return baseMsg
}

func NewValidationError(message string, internalErr ...error) *AppError {
	var details any
	var fieldErrors []FieldError
	
	// Extract validation errors if provided
	if len(internalErr) > 0 && internalErr[0] != nil {
		fieldErrors = ExtractValidationErrors(internalErr[0])
	}
	
	// Only include details if we have field errors
	if len(fieldErrors) > 0 {
		details = fieldErrors
	} else {
		details = nil // omitempty will hide this in JSON
	}
	
	ae := &AppError{
		Status:  http.StatusBadRequest,
		Code:    "VALIDATION_ERROR",
		Message: message,
		Details: details,
	}
	
	if len(internalErr) > 0 {
		ae.Err = internalErr[0]
	}
	
	return ae
}

func ExtractValidationErrors(err error) []FieldError {
	var fieldErrors []FieldError
	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) { 
		for _, fieldErr := range validationErrors {
			fieldErrors = append(fieldErrors, FieldError{
				Field:   fieldErr.Field(), 
				Message: getValidationErrorMessage(fieldErr),
			})
		}
	} else {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   "general",
			Message: err.Error(), 
		})
	}
	return fieldErrors
}


func getValidationErrorMessage(fieldErr validator.FieldError) string {
	switch fieldErr.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldErr.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fieldErr.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", fieldErr.Field(), fieldErr.Param())
	case "eqfield":
		return fmt.Sprintf("%s does not match %s", fieldErr.Field(), fieldErr.Param())
	default:
		return fmt.Sprintf("%s is invalid", fieldErr.Field())
	}
}