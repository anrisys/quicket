package models

import (
	"time"

	"gorm.io/gorm"
)

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
