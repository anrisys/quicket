package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EventSnapshot struct {
	ID       uint64 `gorm:"primaryKey"`
	PublicID string `gorm:"type:char(36);not null;uniqueIndex:idx_events_public_id"`

	Title       string  `gorm:"type:varchar(256);not null"`
	Category    string  `gorm:"type:varchar(100);not null;index:idx_events_category"`
	Description *string `gorm:"type:text"`

	LocationCity    string `gorm:"type:varchar(100);not null;index:idx_events_city"`
	LocationCountry string `gorm:"type:varchar(100);not null"`

	StartDate time.Time `gorm:"type:datetime;not null;index:idx_events_start_date"`
	EndDate   time.Time `gorm:"type:datetime;not null"`

	BasePrice      float64 `gorm:"type:decimal(10,2);not null"`
	MaxSeats       uint64  `gorm:"not null"`
	AvailableSeats uint64  `gorm:"not null"`

	Attributes datatypes.JSON `gorm:"type:json"`

	CreatedAt time.Time `gorm:"type:datetime(3);not null"`
	UpdatedAt time.Time `gorm:"type:datetime(3);not null"`

	DeletedAt gorm.DeletedAt `gorm:"type:datetime(3);index:idx_events_deleted_at"`
}
