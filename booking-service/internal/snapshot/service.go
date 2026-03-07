package snapshot

import (
	"context"

	"github.com/rs/zerolog"
)

/*
|--------------------------------------------------------------------------
| Event Snapshot Service
|--------------------------------------------------------------------------
*/
type EventSnapshotService struct {
	esr EventSnapshotRepository
	l   zerolog.Logger
}

func NewEventSnapshotService(esr EventSnapshotRepository, l zerolog.Logger) *EventSnapshotService {
	return &EventSnapshotService{esr: esr, l: l}
}

func (ess *EventSnapshotService) FindIDsByPublicID(ctx context.Context, publicID string) (*EventIDs, error) {
	return ess.esr.FindIDsByPublicID(ctx, publicID)
}

func (ess *EventSnapshotService) FindSeatBasePrice(ctx context.Context, id uint64) (float64, error) {
	return ess.esr.FindEventSeatPrice(ctx, id)
}

/*
|--------------------------------------------------------------------------
| User Snapshot Service
|--------------------------------------------------------------------------
*/
type UserSnapshotService struct {
	usr UserSnapshotRepository
	l   zerolog.Logger
}

func NewUserSnapshotService(usr UserSnapshotRepository, l zerolog.Logger) *UserSnapshotService {
	return &UserSnapshotService{usr: usr, l: l}
}

func (ess *UserSnapshotService) GetUserPrimaryID(ctx context.Context, publicID string) (uint64, error) {
	return ess.usr.FindUserPrimaryIDByPublicID(ctx, publicID)
}
