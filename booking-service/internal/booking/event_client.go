package booking

import "context"

type EventClientService interface {
	ReserveSeats(ctx context.Context, eventPublicID string, seats uint64) (float64, error)
	ReleaseSeats(ctx context.Context, eventPublicID string, seats uint64) error
}
