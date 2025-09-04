package dto

import "time"

type CreateEventRequest struct {
	Title     string 	`json:"title" binding:"required,min=3,max=256"`
	StartDate time.Time `json:"start_date" binding:"required,gttoday"`
	EndDate time.Time `json:"end_date" binding:"required,gtefield=StartDate"`
	Description string `json:"description" binding:"max=2000,omitempty"`
	MaxSeats uint64 `json:"max_seats" binding:"required,gt=0"`
}