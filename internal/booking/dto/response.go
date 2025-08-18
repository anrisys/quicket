package dto

type ResponseSuccess struct {
	Code    string `json:"code" example:"SUCCESS"`
	Message string `json:"message" example:"Operation successful"`
}

type CreateBookingSuccessResponse struct {
	ResponseSuccess `json:",inline"`
	Booking         BookingDTO `json:",inline"`
}