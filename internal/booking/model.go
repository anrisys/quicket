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
	Status string `gorm:"column:status;type:ENUM('success', 'failed', 'pending');default:'pending'"`
	ExpiredAt time.Time `gorm:"column:expired_at;not null"`
}

func (b *Booking) TableName() string {
	return "bookings"
}