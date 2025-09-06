package internal

import "errors"

var (
	ErrEventAlreadyExist = errors.New("event already exists")
	ErrEventNotFound = errors.New("event not found")
	ErrDB = errors.New("database error")
)