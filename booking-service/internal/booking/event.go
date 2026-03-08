package booking

type BookingCreatedEvent struct {
	BookingPublicID string
	EventID         uint64
	UserID          uint64
	Seats           uint64
}

type BookingCancelledEvent struct {
	BookingPublicID string
}

type BookingReleaseSeats struct {
	EventPublicID string
	Seats         uint64
}
