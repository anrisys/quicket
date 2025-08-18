package dto

import "time"

type EventDTO struct {
	ID             	uint
	PublicID       	string   
	Title          	string   
	Description    	*string  
	StartDate      	time.Time
	EndDate        	time.Time
	MaxSeats       	uint64   
	AvailableSeats 	uint64
	CreatedAt		time.Time
	UpdatedAt		time.Time
}

type SimpleEventDTO struct {
	PublicID       	string   	`json:"public_id" example:"evt_123"`
	Title          	string		`json:"title" example:"Concert Night"`
	StartDate      	time.Time	`json:"start_date" example:"2023-12-31T20:00:00Z"`
	EndDate        	time.Time 	`json:"end_date" example:"2023-12-31T23:59:59Z"`
}