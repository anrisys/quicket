package clients

import "errors"

var (
	ErrServiceClientUserNotFound = errors.New("user not found")
	ErrRequestClientFailed = errors.New("http request failed")
	ErrClientUnexpectedResponseStatus = errors.New("unexpected response status")
	ErrClientFailedToDecodeResponse = errors.New("failed to decode response")
)