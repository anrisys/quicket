package payment

import (
	"context"
	"fmt"
	"quicket-payment-service/internal/domain"
	"quicket-payment-service/internal/payment/dto"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Usecase struct {
	r   Repository
	pg  PaymentGateway
	evp EventPublisher
	l   zerolog.Logger
}

func NewUsecase(r Repository, pg PaymentGateway, evp EventPublisher, l zerolog.Logger) *Usecase {
	return &Usecase{r: r, pg: pg, evp: evp, l: l}
}

func (u *Usecase) InitPayment(ctx context.Context, data dto.InitPaymentRequest) error {
	// Check if booking.created event already expired or not.
	var status domain.PaymentStatus
	if data.ExpiredAt.Before(time.Now()) {
		status = domain.PaymentStatusExpired
	} else {
		status = domain.PaymentStatusPending
	}

	// Create payment session
	py := PaymentSession{
		ID:        uuid.NewString(),
		BookingID: data.BookingID,
		Amount:    data.Amount,
		Currency:  data.Currency,
		Status:    status,
	}

	saved, err := u.r.CreateSession(ctx, &py)
	if err != nil {
		return fmt.Errorf("usecase.InitPayment: %w", err)
	}

	// Hit Payment Gateway API
	req := dto.CreatePaymentRequest{
		ExternalID: saved.ID,
		Amount:     saved.Amount,
		Currency:   saved.Currency,
		Expiry:     saved.ExpiredAt,
	}
	_, err = u.pg.CreatePayment(ctx, req)

	// Publish event payment initilized
	err = u.evp.PaymentInitialized(ctx, *saved)

	if err != nil {
		return fmt.Errorf("usecase.InitPayment: %w", err)
	}

	return nil
}
