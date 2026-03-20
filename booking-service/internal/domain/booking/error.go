package booking

import "errors"

var (
	ErrBookingNotFound = errors.New("booking not found")

	ErrSeatsNotAvailable = errors.New("seats not available")

	ErrEventNotFound           = errors.New("event not found")
	ErrEventClosed             = errors.New("event closed")
	ErrEventServiceUnavailable = errors.New("event service unavailable")

	ErrInvalidStatusTransition = errors.New("invalid booking status transition")
)
