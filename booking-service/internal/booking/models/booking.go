package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type BookingStatus string
type PaymentMethod string
type Channel string

const (
	// Booking Status
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusSuccess   BookingStatus = "success"
	BookingStatusFailed    BookingStatus = "failed"
	BookingStatusCancelled BookingStatus = "cancelled"

	// Payment Method
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
	PaymentMethodEWallet      PaymentMethod = "e_wallet"
	PaymentMethodCash         PaymentMethod = "cash"

	// Channel
	ChannelWeb     Channel = "web"
	ChannelMobile  Channel = "mobile"
	ChannelPartner Channel = "partner"
)

type Booking struct {
	ID       uint64 `gorm:"primaryKey;autoIncrement"`
	PublicID string `gorm:"type:char(36);not null;uniqueIndex:idx_bookings_public_id"`

	EventID uint64 `gorm:"not null;index:idx_bookings_event_id"`
	UserID  uint64 `gorm:"not null;index:idx_bookings_user_id"`

	Seats      uint64  `gorm:"not null"`
	TotalPrice float64 `gorm:"type:decimal(12,2);not null"`
	Currency   string  `gorm:"type:char(3);not null;default:USD"`

	Status BookingStatus `gorm:"type:enum('pending','success','failed','cancelled');not null;default:'pending';index:idx_bookings_status'"`

	PaymentMethod PaymentMethod `gorm:"type:enum('credit_card','bank_transfer','e_wallet','cash');not null;index:idx_bookings_payment_method'"`

	Channel Channel `gorm:"type:enum('web','mobile','partner');not null;index:idx_bookings_channel'"`

	ConfirmedAt *time.Time `gorm:"type:datetime(3)"`
	ExpiredAt   time.Time  `gorm:"type:datetime(3);not null"`

	Metadata datatypes.JSON `gorm:"type:json"`

	CreatedAt time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3);index:idx_bookings_created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3)"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime(3);index:idx_bookings_deleted_at"`
}
