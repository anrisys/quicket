package booking

import "context"

type EventClientService interface {
	ReserveSeats(ctx context.Context, eventPublicID string, seats uint64) (uint64, error)
}
