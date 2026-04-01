package infrastructure

type clientServiceError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type reserveSeatsSuccessResponse struct {
	SeatPrice uint64 `json:"seat_price"`
}
