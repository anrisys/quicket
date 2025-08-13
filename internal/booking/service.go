package booking

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anrisys/quicket/internal/booking/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/anrisys/quicket/pkg/util"
	"github.com/rs/zerolog"
)

type UserReader interface {
	GetUserID(ctx context.Context, publicID string) (*int, error)
}

type EventReader interface {
	GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*dto.EventDateTimeAndSeats, error)
}

type Service interface {
	Create(ctx context.Context, req *dto.CreateBookingRequest, userID string) (*dto.BookingDTO, error)
}

type service struct {
	repo Repository
	events EventReader
	users UserReader
	logger zerolog.Logger
}

func NewService(repo Repository, 
	events EventReader, 
	logger zerolog.Logger,
	users UserReader) *service {
	return &service{
		repo: repo,
		events: events,
		users: users,
		logger: logger,
	}
}

func (s *service) Create(ctx context.Context, req *dto.CreateBookingRequest, userPublicID string) (*dto.BookingDTO, error) {
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

	log.Info().
		Str("booking_public_id", persisted.PublicID).
		Uint("event_id", persisted.EventID).
		Uint("user_id", persisted.UserID).
		Msg("Booking created")

	bDTO := s.prepareBookingDTO(ctx, persisted, req.EventID, userPublicID)
	return bDTO, nil
}

func (s *service) prepareBooking(ctx context.Context, eventID uint, userID int, seats uint) (*Booking, error) {
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

func (s *service) prepareBookingDTO(_ctx context.Context, booking *Booking, eventPublicID string, userPublicID string) *dto.BookingDTO {
	return &dto.BookingDTO{
		PublicID: booking.PublicID,
		EventID: eventPublicID,
		UserID: userPublicID,
		Seats: booking.Seats,
		Status: booking.Status,
	}
}