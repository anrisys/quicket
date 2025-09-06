package internal

import "time"

type EventDTO struct {
	ID             uint
	PublicID       string
	Title          string
	Description    *string
	StartDate      time.Time
	EndDate        time.Time
	MaxSeats       uint64
	AvailableSeats uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type SimpleEventDTO struct {
	PublicID  string    `json:"public_id" example:"evt_123"`
	Title     string    `json:"title" example:"Concert Night"`
	StartDate time.Time `json:"start_date" example:"2023-12-31T20:00:00Z"`
	EndDate   time.Time `json:"end_date" example:"2023-12-31T23:59:59Z"`
}

type CreateEventRequest struct {
	Title     string 	`json:"title" binding:"required,min=3,max=256"`
	StartDate time.Time `json:"start_date" binding:"required,gttoday"`
	EndDate time.Time `json:"end_date" binding:"required,gtefield=StartDate"`
	Description string `json:"description" binding:"max=2000,omitempty"`
	MaxSeats uint64 `json:"max_seats" binding:"required,gt=0"`
}

type ResponseSuccess struct {
	Code    string `json:"code" example:"SUCCESS"`
	Message string `json:"message" example:"Operation successful"`
}

type CreateEventSuccessResponse struct {
	ResponseSuccess `json:",inline"`
	Event           SimpleEventDTO `json:"event"`
}

type EventDateTimeAndSeats struct {
	ID             int
	AvailableSeats uint64
	EndDate        time.Time
}

type UserDTO struct {
	ID       int
	Email    string
	PublicID string
	Role     string
}