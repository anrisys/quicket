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