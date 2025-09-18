package usersnapshot

import "errors"

var (
	ErrUserNotFound = errors.New("user snapshot not found")
	ErrDB = errors.New("database error")
)