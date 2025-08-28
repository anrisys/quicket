package types

import (
	"context"

	commonDTO "github.com/anrisys/quicket/internal/dto"
)

type UserReader interface {
	GetUserID(ctx context.Context, publicID string) (*uint, error)
	FindUserByPublicID(ctx context.Context, publicID string) (*commonDTO.UserDTO, error)
}

type EventReader interface {
	GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*commonDTO.EventDateTimeAndSeats, error)
}

type BookingReader interface {
	GetSimpleBookingDTO(ctx context.Context, publicID string) (*commonDTO.SimpleBookingDTO, error)
}

type SimulatePayment interface {
	SimulatePayment(ctx context.Context, bookData *commonDTO.SimulateBookingPayment) (*commonDTO.PaymentDTO, error)
}