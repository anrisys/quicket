package errs

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err 	error `json:"-"`
	Metadata interface{} `json:"metadata,omitempty"`
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

var (
	ErrInternal = &AppError{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_ERROR",
		Message: "An unexpected internal server error occurred. Please try again later.",
	}

	ErrNotFound = &AppError{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: "The requested resource was not found.",
	}

	ErrBadRequest = &AppError{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: "The request could not be understood or was invalid.",
	}

	ErrUnauthorized = &AppError{
		Status:  http.StatusUnauthorized,
		Code:    "UNAUTHORIZED",
		Message: "Authentication required or invalid credentials.",
	}

	ErrForbidden = &AppError{
		Status:  http.StatusForbidden,
		Code:    "FORBIDDEN",
		Message: "You do not have permission to access this resource.",
	}
)