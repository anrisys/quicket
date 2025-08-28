package dto

type CreateBookingRequest struct {
	Seats uint `json:"seats" binding:"required,gt=0"`
}