package booking

/* LEGACY ERROR
var (
	ErrBookingNotFound = errors.New("booking not found")
	ErrUserNotFound = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")
	ErrSeatsUnavailable = errors.New("no available seats")
	ErrNotEnoughSeats = errors.New("not enough setas")
	ErrDB = errors.New("database error")
)
*/

const (
	CodeValidation   = "VALIDATION_ERROR"
	CodeNotFound     = "NOT_FOUND"
	CodeConflict     = "CONFLICT"
	CodeInternal     = "INTERNAL_ERROR"
	CodeUnauthorized = "UNAUTHORIZED"
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
