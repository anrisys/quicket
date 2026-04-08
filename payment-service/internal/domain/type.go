package domain

type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusFailed     PaymentStatus = "failed"
	PaymentStatusSuccess    PaymentStatus = "success"
	PaymentStatusExpired    PaymentStatus = "expired"
	PaymentStatusRefunded   PaymentStatus = "refunded"
)
