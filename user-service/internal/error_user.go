package internal

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrDB = errors.New("database error")
)