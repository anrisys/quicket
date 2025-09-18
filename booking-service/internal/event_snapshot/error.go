package eventsnapshot

import "errors"

var (
	ErrBookingNotFound  = errors.New("booking not found")
	ErrEventNotFound    = errors.New("event not found")
	ErrDB               = errors.New("database error")
)