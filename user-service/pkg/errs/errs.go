package errs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err 	error `json:"-"`
	Details any `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("AppError (Code: %s, Message: %s, Status: %d) -> Wrapped Error: %v",
			e.Code, e.Message, e.Status, e.Err)
	}
	return fmt.Sprintf("AppError (Code: %s, Message: %s, Status: %d)", e.Code, e.Message, e.Status)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool  {
	if target, ok := target.(*AppError); ok {
		return e.Code == target.Code
	}
	return false
}

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

func NewAppError(status int, code, message string, internalErr ...error) *AppError {
	ae := &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: nil,
	}
	if len(internalErr) > 0 {
		ae.Err = internalErr[0] // Only take the first error if multiple are passed
	}
	return ae
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

func NewConflictError(message string, internalErr ...error) *AppError {
	ae := &AppError{
		Status:  http.StatusConflict,
		Code:    "CONFLICT_ERROR",
		Message: message,
		Details: nil,
	}
	if len(internalErr) > 0 {
		ae.Err = internalErr[0]
	}
	return ae
}

func NewServiceUnavailableError(message string, internalErr ...error) *AppError {
    ae := &AppError{
		Status:  http.StatusServiceUnavailable,
		Code:    "SERVICE_UNAVAILABLE",
		Message: message,
		Details: nil,
	}
	if len(internalErr) > 0 {
		ae.Err = internalErr[0]
	}
	return ae
}

func NewInternalError(message string, internalErr ...error) *AppError {
	ae := &AppError{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_ERROR",
		Message: message,
		Details: nil,
	}
	if len(internalErr) > 0 {
		ae.Err = internalErr[0]
	}
	return ae
}

func NewErrNotFound(resource string, internalErr ...error) *AppError {
    	ae := &AppError{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
		Details: nil,
	}
	if len(internalErr) > 0 {
		ae.Err = internalErr[0]
	}
	return ae
}

var (
	ErrInternal = NewAppError(
		http.StatusInternalServerError, 
		"INTERNAL_ERROR", 
		"An unexpected internal server error occurred. Please try again later.",
	)
	ErrNotFound = NewAppError(
		http.StatusNotFound, 
		"NOT_FOUND", 
		"The requested resource was not found.",
	)
	ErrBadRequest = NewAppError(
		http.StatusBadRequest, 
		"BAD_REQUEST", 
		"The request could not be understood or was invalid.")
	ErrUnauthorized = NewAppError(
		http.StatusUnauthorized, 
		"UNAUTHORIZED", 
		"Authentication required or invalid credentials.",
	)
	ErrForbidden = NewAppError(
		http.StatusForbidden, 
		"FORBIDDEN", 
		"You do not have permission to access this resource.",
	)
	ErrConflict = NewAppError(
		http.StatusConflict,
		"CONFLICT",
		"The request could not be completed due to a conflict with the current state of the resource.",
	)
	ErrServiceUnavailable = NewAppError(
		http.StatusServiceUnavailable,
		"SERVICE_UNAVAILABLE",
		"The service is temporarily unavailable. Please try again later.",
	)
	ErrTooManyRequests = NewAppError(
		http.StatusTooManyRequests,
		"TOO_MANY_REQUESTS",
		"You have sent too many requests in a given amount of time. Please try again later.",
	)
)

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