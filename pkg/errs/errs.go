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
	}
	if len(internalErr) > 0 {
		ae.Err = internalErr[0] // Only take the first error if multiple are passed
	}
	return ae
}

func NewValidationError(message string, fields []FieldError, internalErr ...error) *ValidationError {
	ve := &ValidationError{
		AppError: AppError{
			Status:  http.StatusBadRequest, // Validation errors are typically Bad Request
			Code:    "VALIDATION_ERROR",
			Message: message,
		},
		Fields: fields,
	}
	if len(internalErr) > 0 {
		ve.Err = internalErr[0]
	}
	return ve
}

func NewConflictError(message string) *AppError {
    return NewAppError(http.StatusConflict, "CONFLICT", message, nil)
}

func NewServiceUnavailableError(message string) *AppError {
    return NewAppError(http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message, nil) 
}

func NewInternalError(message string) *AppError {
    return NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", message, nil)
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
