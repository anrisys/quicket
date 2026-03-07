package dto

import (
	"quicket/booking-service/internal/domain/booking"
	"time"
)

type DateRangeFilterOption struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

type BookingFilterOptions struct {
	Status        []booking.BookingStatus `json:"status,omitempty"`
	PaymentMethod []booking.PaymentMethod `json:"payment_method,omitempty"`
	Channel       []booking.Channel       `json:"channel,omitempty"`
	CreatedAt     *DateRangeFilterOption  `json:"created_at,omitempty"`
	ConfirmedAt   *DateRangeFilterOption  `json:"confirmed_at,omitempty"`
}

type EventFilterOptions struct {
	EventPublicID  string                 `json:"event_public_id,omitempty"`
	EventCategory  string                 `json:"event_category,omitempty"`
	EventStartDate *DateRangeFilterOption `json:"event_start_date,omitempty"`
}

type UserFilterOption struct {
	UserPublicID string `json:"user_public_id,omitempty"`
}

type AdminBookingListFilter struct {
	BookingFilter *BookingFilterOptions `json:"booking,omitempty"`
	EventFilter   *EventFilterOptions   `json:"event,omitempty"`
	UserFilter    *UserFilterOption     `json:"user,omitempty"`
}

type RequestSort struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type PaginationRequest struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset,omitempty"`
}

type AdminBookingListRequest struct {
	Filters    *AdminBookingListFilter `json:"filters,omitempty"`
	Sort       *[]RequestSort          `json:"sort,omitempty"`
	Pagination *PaginationRequest      `json:"pagination,omitempty"`
}

type CreateBookingRequest struct {
	EventPublicID string `json:"event_public_id" binding:"required"`
	Seats         uint64 `json:"seats" binding:"required,min=1"`
	PaymentMethod string `json:"payment_method" binding:"required"`
	Channel       string `json:"channel" binding:"required"`
}
