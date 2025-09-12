package internal

import "errors"

var (
	ErrBookingNotFound = errors.New("booking not found")
	ErrUserNotFound = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")
	ErrSeatsUnavailable = errors.New("no available seats")
	ErrNotEnoughSeats = errors.New("not enough setas")
	ErrDB = errors.New("database error")
)