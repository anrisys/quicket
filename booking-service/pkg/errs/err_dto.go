package errs

import "fmt"

type ErrorResponse struct {
	Code    string       `json:"code" example:"VALIDATION_ERROR"`
	Message string       `json:"message" example:"Invalid input data"`
	Fields  []FieldError `json:"fields,omitempty"`
}

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
	Details any    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("AppError (Code: %s, Message: %s, Status: %d) -> Wrapped Error: %v",
			e.Code, e.Message, e.Status, e.Err)
	}
	return fmt.Sprintf("AppError (Code: %s, Message: %s, Status: %d)", e.Code, e.Message, e.Status)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) Is(target error) bool {
	if target, ok := target.(*AppError); ok {
		return e.Code == target.Code
	}
	return false
}
