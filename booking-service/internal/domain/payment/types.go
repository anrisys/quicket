package payment

type PaymentStatus string

const (
	PaymentStatusFailed  PaymentStatus = "failed"
	PaymentStatusSuccess PaymentStatus = "success"
)
