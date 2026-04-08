package payment

import "context"

type Repository interface {
	CreateSession(ctx context.Context, data *PaymentSession) (*PaymentSession, error)
}
