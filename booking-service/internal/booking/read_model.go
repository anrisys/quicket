package booking

import (
	"quicket/booking-service/internal/domain/booking"
	"time"
)

type BookingRow struct {
	BookingPublicID    string                `gorm:"column:booking_public_id"`
	BookingStatus      booking.BookingStatus `gorm:"column:booking_status"`
	BookingSeats       uint64                `gorm:"column:booking_seats"`
	BookingTotalPrice  float64               `gorm:"column:booking_total_price"`
	BookingCurrency    string                `gorm:"column:booking_currency"`
	BookingChannel     booking.Channel       `gorm:"column:booking_channel"`
	BookingCreatedAt   time.Time             `gorm:"column:booking_created_at"`
	BookingConfirmedAt *time.Time            `gorm:"column:booking_confirmed_at"`

	EventPublicID     string    `gorm:"column:event_public_id"`
	EventTitle        string    `gorm:"column:event_title"`
	EventCategory     string    `gorm:"column:event_category"`
	EventStartDate    time.Time `gorm:"event_start_date"`
	EventLocationCity string    `gorm:"event_location_city"`

	UserPublicID string `gorm:"column:user_public_id"`
	UserFullName string `gorm:"column:user_full_name"`
	UserEmail    string `gorm:"column:user_email"`
}

type ExpiredBookingRow struct {
	ID            uint64 `gorm:"column:ID"`
	EventPublicID string `gorm:"event_public_id"`
	Seats         uint64 `gorm:"seats"`
}

type BookingCancellationRow struct {
	ID              uint64                `gorm:"column:ID"`
	BookingPublicID string                `gorm:"column:booking_public_id"`
	BookingStatus   booking.BookingStatus `gorm:"column:booking_status"`
	BookingSeats    uint64                `gorm:"column:booking_seats"`

	EventPublicID string `gorm:"column:event_public_id"`

	UserPublicID string `gorm:"column:user_public_id"`
}
