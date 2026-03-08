package consumer

import (
	"time"
)

type EventCreatedMessage struct {
	ID        		uint 			`json:"id"`
	PublicID 		string 			`json:"public_id"`
	Title 			string 			`json:"title"`
	StartDate 		time.Time 		`json:"start_date"`
	EndDate 		time.Time 		`json:"end_date"`
	AvailableSeats 	uint64 			`json:"available_seats"`
	CreatedAt 		time.Time		`json:"created_at"`
	UpdatedAt 		time.Time		`json:"updated_at"`
	Version       	uint       		`json:"version"`
}

type EventUpdatedMessage struct {
	EventCreatedMessage `json:",inline"`
}

type UserCreatedMessage struct {
	ID 			uint 	`json:"id"`
	PublicID 	string 	`json:"public_id"`
}

type UserDeletedMessage struct {
	ID uint `json:"id"`
}