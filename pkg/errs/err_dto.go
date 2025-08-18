package errs

type ErrorResponse struct {
	Code    string       `json:"code" example:"VALIDATION_ERROR"`
	Message string       `json:"message" example:"Invalid input data"`
	Fields  []FieldError `json:"fields,omitempty"`
}