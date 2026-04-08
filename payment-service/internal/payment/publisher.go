package payment

import "context"

type EventPublisher interface {
	PaymentInitialized(ctx context.Context, payload PaymentSession) error
}
