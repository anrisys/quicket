package payment

import (
	"context"
	"quicket-payment-service/internal/payment/dto"
)

type PaymentGateway interface {
	CreatePayment(ctx context.Context, req dto.CreatePaymentRequest) (*dto.CreatePaymentResponse, error)
}
