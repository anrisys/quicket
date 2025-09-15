package mq

import "errors"

var (
	ErrFailedToDeclareExchange = errors.New("failed to declare exchange")
	ErrFailedToDeclareQueue = errors.New("failed to declare queue")
)