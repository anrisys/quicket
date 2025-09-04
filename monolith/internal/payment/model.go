package payment

import "time"

type Payment struct {
	ID uint `gorm:"primarykey"`
	PublicID 	string `gorm:"column:public_id;type:char(36);uniqueIndex;not null"`
	Amount float32 `gorm:"column:amount;not null"`
	Status string `gorm:"column:status;type:ENUM('success', 'failed');not null"`
	BookingID uint `gorm:"column:booking_id;not null"`
	UserID uint `gorm:"column:user_id;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Payment) TableName() string {
	return "payments"
}