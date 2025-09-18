package eventsnapshot

import "time"

type EventDateTimeAndSeats struct {
	ID             int
	AvailableSeats uint64
	EndDate        time.Time
}