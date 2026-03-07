package infrastructure

type clientServiceError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type reserveSeatsSuccessResponse struct {
	SeatPrice float64 `json:"seat_price"`
}
