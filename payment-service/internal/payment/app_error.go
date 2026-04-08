package payment

const (
	CodeValidation   = "VALIDATION_ERROR"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
	CodeForbidden    = "FORBIDDEN"
	CodeBadRequest   = "BAD_REQUEST"
)

type AppError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, code, message string, err error) *AppError {
	return &AppError{
		Status:  status,
		Code:    code,
		Message: message,
		Err:     err,
	}
}
