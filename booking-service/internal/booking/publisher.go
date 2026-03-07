package booking

import "context"

type EventPublisher interface {
	PublishBookingCreated(ctx context.Context, payload BookingCreatedEvent) error
	PublishBookingCancelled(ctx context.Context, payload BookingCancelledEvent) error
}
