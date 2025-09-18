package eventsnapshot

import (
	"time"
)

type EventSnapshot struct {
	ID        		uint 		`gorm:"primarykey"`
	PublicID 		string 		`gorm:"column:public_id;type:char(36);uniqueIndex"`
	Title 			string 		`gorm:"column:title;size:256;not null"`
	StartDate 		time.Time 	`gorm:"column:start_date;not null;index"`
	EndDate 		time.Time 	`gorm:"column:end_date;not null"`
	AvailableSeats 	uint64 		`gorm:"column:available_seats"`
	UpdatedAt 		time.Time
	Version			uint		`gorm:"column:version"`
}

func (e *EventSnapshot) TableName() string {
	return "events_snapshot"
}