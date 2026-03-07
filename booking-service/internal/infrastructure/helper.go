package infrastructure

import (
	"context"
	"quicket/booking-service/internal/domain/booking"
	"time"
)

func Retry(
	ctx context.Context,
	attempts int,
	initialDelay time.Duration,
	fn func() error,
) error {

	delay := initialDelay
	var err error

	for i := 0; i < attempts; i++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		err = fn()
		if err == nil {
			return nil
		}

		// last attempt → return error
		if i == attempts-1 {
			break
		}

		select {
		case <-time.After(delay):
			delay *= 2
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return err
}

func mapEventServiceError(code string) error {
	switch code {

	case "EVENT_NOT_FOUND":
		return booking.ErrEventNotFound

	case "SEATS_NOT_AVAILABLE":
		return booking.ErrSeatsNotAvailable

	case "EVENT_CANCELLED", "EVENT_ALREADY_ENDED":
		return booking.ErrEventClosed

	default:
		return booking.ErrEventServiceUnavailable
	}
}
