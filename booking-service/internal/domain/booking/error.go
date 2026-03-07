package booking

import "errors"

var (
	ErrSeatsNotAvailable       = errors.New("seats not available")
	ErrBookingNotFound         = errors.New("booking not found")
	ErrEventNotFound           = errors.New("event not found")
	ErrInvalidStatusTransition = errors.New("invalid booking status transition")
	ErrEventClosed             = errors.New("event closed")
	ErrEventServiceUnavailable = errors.New("event service unavailable")
)
