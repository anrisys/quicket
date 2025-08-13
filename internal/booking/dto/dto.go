package dto

import "time"

type BookingDTO struct {
	PublicID string `json:"id"`
	EventID  string `json:"event_id"`
	UserID   string `json:"user_id"`
	Seats    uint   `json:"number_of_seats"`
	Status   string `json:"status"`
}

type EventDateTimeAndSeats struct {
	ID             	int
	AvailableSeats 	uint64
	EndDate 		time.Time
}