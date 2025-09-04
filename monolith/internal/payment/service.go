package payment

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	commonDTO "github.com/anrisys/quicket/internal/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/anrisys/quicket/pkg/util"
	"github.com/rs/zerolog"
)

const (
	numWorkers = 5
	jobQueueSize = 100
)

type PaymentJob struct {
	BookingID uint
	UserID uint
	Amount float32
	PublicID string
}

type PaymentServiceInterface any

type PaymentService struct {
	r *GormRepository
	logger zerolog.Logger
	jobQueue chan PaymentJob
}

func NewPaymentService(r *GormRepository, logger zerolog.Logger) *PaymentService {
	jobQueue := make(chan PaymentJob, jobQueueSize)
	for i := 1; i <= numWorkers; i++ {
		go startWorker(i, r, logger, jobQueue)
	}
	return &PaymentService{
		r: r,
		logger: logger,
		jobQueue: jobQueue,
	}
}

func (s *PaymentService) SimulatePayment(ctx context.Context, bookData *commonDTO.SimulateBookingPayment) (*commonDTO.PaymentDTO, error) {
	publicID, err := util.GeneratePublicID(ctx)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to generate payment public ID")
		return nil, fmt.Errorf("payment#SimulatePayment: %w", err)
	}
	
	job := PaymentJob{
		BookingID: bookData.BookingID,
		UserID: bookData.UserID,
		Amount: bookData.Amount,
		PublicID: publicID,
	}

	select {
		case s.jobQueue <- job:
			s.logger.Info().
				Str("payment_public_id", publicID).
				Uint("booking_id", bookData.BookingID).
				Int("queue_length", len(s.jobQueue)).
				Msg("payment simulation added into queue")
		default: 
			s.logger.Warn().
				Str("payment_public_id", publicID).
				Msg("job queue is full, payment simulation rejected")
			return nil, errs.NewServiceUnavailableError("payment system busy")
	}

	return nil, nil
}

func startWorker(id int, r *GormRepository, logger zerolog.Logger, jobQueue <-chan PaymentJob) {
	logger.Info().Int("worker_id", id).Msg("payment worker started")

	defer func() {
		if r := recover(); r != nil {
			logger.Error().
				Int("worker_id", id).
				Msg("worker recovered from panice")
		}
	}()
	for job := range jobQueue {
		logger.Info().
			Int("worker_id", id).
			Str("payment_public_id", job.PublicID).
			Uint("booking_id", job.BookingID).
			Msg("start process payment")
		
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

		status := "failed"
		if rand.Intn(10) < 8 {
			status = "success"
		}

		ctx := context.Background()
		p := &Payment{
			PublicID: job.PublicID,
			Amount: job.Amount,
			Status: status,
			BookingID: job.BookingID,
			UserID: job.UserID,
		}
		if _, err := r.CreatePaymentAndUpdateBookingStatus(ctx, p); err != nil {
			logger.Error().Err(err).
				Int("worker_id", id).
				Str("public_payment_id", job.PublicID).
				Uint("booking_id", p.BookingID).
				Uint("user_id", p.UserID).
				Msg("failed to update and create payment")
				continue
		}

		logger.Info().
			Int("worker_id", id).
			Uint("booking_id", job.BookingID).
			Str("status", status).
			Msg("Payment job completed")			
	}
}