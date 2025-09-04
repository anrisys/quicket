package dto

import "time"

type BookingDTO struct {
	PublicID  string    `json:"id"`
	EventID   string    `json:"event_id"`
	UserID    string    `json:"user_id"`
	Seats     uint      `json:"seats"`
	Status    string    `json:"status"`
	ExpiredAt time.Time `json:"expired_at"`
}