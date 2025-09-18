package booking

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	PublicID string `gorm:"column:public_id;type:char(36);uniqueIndex"`
	EventID uint `gorm:"column:event_id;not null"`
	UserID uint `gorm:"column:user_id;not null"`
	Seats uint `gorm:"column:seats;not null"`
	TotalPrice float32 `gorm:"column:total_price;not null"`
	Status string `gorm:"column:status;type:ENUM('success', 'failed', 'pending');default:'pending'"`
	ExpiredAt time.Time `gorm:"column:expired_at;not null"`
}

func (b *Booking) TableName() string {
	return "bookings"
}

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

type User struct {
	gorm.Model
	PublicID string `gorm:"column:public_id;type:char(36);uniqueIndex"`
}