package errs

import (
	"fmt"
	"net/http"
)

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