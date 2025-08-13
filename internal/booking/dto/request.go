package dto

type CreateBookingRequest struct {
	EventID string `json:"event_id" binding:"required"`
	Seats   uint   `json:"seats" binding:"required;gt=0"`
}