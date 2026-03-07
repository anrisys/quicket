package dto

type Pagination struct {
	Limit        int  `json:"limit"`
	Offset       int  `json:"offset"`
	TotalRecords int  `json:"total_records"`
	TotalPages   int  `json:"total_pages"`
	HasNext      bool `json:"has_next"`
}

type Meta struct {
	Pagination *Pagination `json:"pagination,omitempty"`
}

type SuccessResponse[T any] struct {
	Data T     `json:"data"`
	Meta *Meta `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Fields  interface{} `json:"fields,omitempty"`
}

func NewSuccessResponse[T any](data T) SuccessResponse[T] {
	return SuccessResponse[T]{
		Data: data,
	}
}

func NewPaginatedResponse[T any](data T, pagination Pagination) SuccessResponse[T] {
	return SuccessResponse[T]{
		Data: data,
		Meta: &Meta{
			Pagination: &pagination,
		},
	}
}

type CreateBookingData struct {
	BookingPublicID string `json:"booking_public_id"`
	Status          string `json:"status"`
	ExpiredAt       string `json:"expired_at"`
	TotalPrice      string `json:"total_price"`
	Currency        string `json:"currency"`
}
