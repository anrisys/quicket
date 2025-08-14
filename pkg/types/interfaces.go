package types

import (
	"context"

	commonDTO "github.com/anrisys/quicket/internal/dto"
)

type UserReader interface {
	GetUserID(ctx context.Context, publicID string) (*int, error)
	FindUserByPublicID(ctx context.Context, publicID string) (*commonDTO.UserDTO, error)
}

type EventReader interface {
	GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*commonDTO.EventDateTimeAndSeats, error)
}