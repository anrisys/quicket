package dto

type BookingDTO struct {
	PublicID string `json:"id"`
	EventID  string `json:"event_id"`
	UserID   string `json:"user_id"`
	Seats    uint   `json:"number_of_seats"`
	Status   string `json:"status"`
}