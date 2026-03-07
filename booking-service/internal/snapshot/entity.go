package snapshot

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// EventSnapshot represents a denormalized read model of event-service data.
// It is NOT the source of truth.
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

// UserSnapshot represents a denormalized read model of user-service data.
// It is NOT the source of truth.
type UserSnapshot struct {
	ID       uint64 `gorm:"primaryKey"`
	PublicID string `gorm:"type:char(36);not null;uniqueIndex:idx_users_public_id"`

	Email    *string `gorm:"type:varchar(255)"`
	FullName *string `gorm:"type:varchar(255)"`

	RegisteredAt time.Time `gorm:"type:datetime(3);not null;index:idx_users_registered_at"`

	CreatedAt time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)"`
	UpdatedAt time.Time `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)"`

	DeletedAt gorm.DeletedAt `gorm:"type:datetime(3);index:idx_users_deleted_at"`
}
