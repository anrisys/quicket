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

type RefundRequiredEvent struct {
	PaymentID string
	Amount    uint64
	Currency  string
}

type BookingFailedEvent struct {
	BookingID     string
	FailureReason string
	TotalPrice    string
	Seats         uint64

	UserID string

	EventID string
}
