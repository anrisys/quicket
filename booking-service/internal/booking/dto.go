package booking

import "time"

type CreateBookingRequest struct {
	EventID 	string 	`json:"event_id" binding:"required"`
	Seats 		uint 	`json:"seats" binding:"required,gt=0"`
}

type BookingDTO struct {
	PublicID  string    `json:"id"`
	EventPublicID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Seats     uint      `json:"seats"`
	Status    string    `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}

type EventDateTimeAndSeats struct {
	ID             int
	AvailableSeats uint64
	EndDate        time.Time
}

type ResponseSuccess struct {
	Code    string `json:"code" example:"SUCCESS"`
	Message string `json:"message" example:"Operation successful"`
}

type CreateBookingSuccessResponse struct {
	ResponseSuccess `json:",inline"`
	Data         BookingDTO `json:"booking"`
}