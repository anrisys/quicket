package payment

import "errors"

var (
	ErrBookingNotFound = errors.New("booking not found")
	ErrDB = errors.New("database error")
)