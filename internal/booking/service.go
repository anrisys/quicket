package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

	bookingDTO "github.com/anrisys/quicket/internal/booking/dto"
	commonDTO "github.com/anrisys/quicket/internal/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/anrisys/quicket/pkg/types"
	"github.com/anrisys/quicket/pkg/util"
	"github.com/rs/zerolog"
)

type ServiceInterface interface {
	Create(ctx context.Context, req *bookingDTO.CreateBookingRequest, userID string) (*bookingDTO.BookingDTO, error)
}

type Service struct {
	repo Repository
	events types.EventReader
	users types.UserReader
	payments types.SimulatePayment
	logger zerolog.Logger
}

func NewService(repo Repository, 
	events types.EventReader, 
	logger zerolog.Logger,
	payments types.SimulatePayment,
	users types.UserReader) *Service {
	return &Service{
		repo: repo,
		events: events,
		users: users,
		payments: payments,
		logger: logger,
	}
}

func (s *Service) Create(ctx context.Context, req *bookingDTO.CreateBookingRequest, userPublicID string) (*bookingDTO.BookingDTO, error) {
	log := s.logger.With().
		Str("event_public_id", req.EventID).
		Uint("seats", req.Seats).
		Str("user_public_id", userPublicID).
		Logger()
	
	log.Info().Msg("Create booking request")
	
	log.Debug().Msg("Checking eventID exist or not")
	ev, err := s.events.GetEventDateTimeAndSeats(ctx, req.EventID)
	if err != nil {
		return nil, fmt.Errorf("booking service#create: %w", err)
	}
	now := time.Now()
	if now.After(ev.EndDate) {
		return nil, errs.NewConflictError("can not book past event")
	}

	userID, err := s.users.GetUserID(ctx, userPublicID)
	if err != nil {
		return nil, fmt.Errorf("booking service#create: %w", err)
	}

	newB, err := s.prepareBooking(ctx, uint(ev.ID), *userID, req.Seats)
	if err != nil {
		return nil, err
	}

	persisted, err := s.repo.Create(ctx, newB)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotEnoughSeats):
			return nil, errs.NewConflictError("not enough available seats")
		case errors.Is(err, ErrEventNotFound):
			return nil, errs.NewErrNotFound("event")
		default:
			return nil, fmt.Errorf("booking#create: persist: %w", err)
		}
	}

	// Payment simulation
	bookData := commonDTO.SimulateBookingPayment{
		Amount: persisted.TotalPrice,
		BookingID: persisted.ID,
		UserID: persisted.UserID,
	}

	s.payments.SimulatePayment(ctx, &bookData)
	log.Info().
		Str("booking_public_id", persisted.PublicID).
		Msg("payment simulation triggered asynchronously")

	log.Info().
		Str("booking_public_id", persisted.PublicID).
		Msg("Booking created")

	bDTO := s.prepareBookingDTO(ctx, persisted, req.EventID, userPublicID)
	return bDTO, nil
}

func (s *Service) GetSimpleBookingDTO(ctx context.Context, publicID string) (*commonDTO.SimpleBookingDTO, error) {
	dto, err := s.repo.FindSimpleDTO(ctx, publicID)
	if err != nil {
		if errors.Is(err, ErrBookingNotFound) {
			return nil, errs.NewErrNotFound("booking not found")	
		}
		return nil, errs.ErrInternal
	}
	return dto, nil
}

func (s *Service) prepareBooking(ctx context.Context, eventID uint, userID int, seats uint) (*Booking, error) {
	publicID, err := util.GeneratePublicID(ctx)
	if err != nil {
		return nil, fmt.Errorf("booking#create: generate public id: %w", err)
	}

	return &Booking{
		PublicID: publicID,
		EventID: eventID,
		UserID: uint(userID),
		Seats: seats,
	}, nil
}

func (s *Service) prepareBookingDTO(_ctx context.Context, booking *Booking, eventPublicID string, userPublicID string) *bookingDTO.BookingDTO {
	return &bookingDTO.BookingDTO{
		PublicID: booking.PublicID,
		EventID: eventPublicID,
		UserID: userPublicID,
		Seats: booking.Seats,
		Status: booking.Status,
	}
}