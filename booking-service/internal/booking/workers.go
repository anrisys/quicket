package booking

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

/*
BEST PRACTICE:
1. LOGGING: wrap error found in the usecase/worker layer and let it bubble up. Don't log the error in worker/usecase layer since it causes log noises

2. ERROR WRAPPING CONVENTION:
"<layer>.<method>: <what you were doing>: %w"
e.g. "userRepo.GetByID: scan result: %w"
*/

type Worker struct {
	rr BookingReadRepo
	wr BookingWriteRepository
	ep EventPublisher
	l  zerolog.Logger
}

func NewWoker(rr BookingReadRepo, wr BookingWriteRepository, ep EventPublisher, l zerolog.Logger) *Worker {
	return &Worker{rr: rr, wr: wr, ep: ep, l: l}
}

func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case <-ticker.C:
			w.processExpiredBooking(ctx)
		case <-ctx.Done():
			w.l.Info().Msg("worker stopped")
			return
		}
	}
}

func (w *Worker) processExpiredBooking(ctx context.Context) {
	// Retrieve expired booking
	expiredBookings, err := w.rr.RetrieveExpiredBookings(ctx, 100)
	if err != nil {
		w.l.Error().Err(err).Msg("worker.ProcessExpiredBooking: retrieved expired bookings failed")
		return
	}

	// Do bulk updation
	err = w.wr.UpdateBookingStatusExpiration(ctx, 100)
	if err != nil {
		w.l.Error().Err(err).Msg("worker.ProcessExpiredBooking: update booking expiration failed")
		return
	}

	// publish seats release
	for _, expBooks := range expiredBookings {
		err := w.ep.ReleaseEventSeats(ctx, &BookingReleaseSeats{EventPublicID: expBooks.EventPublicID, Seats: expBooks.Seats})
		if err != nil {
			w.l.Error().Err(err).Msg("worker: processExpiredBooking: failed publish seat release")
		}
	}
}
