package internal

import (
	"context"
	"errors"
	"fmt"
	"quicket/booking-service/pkg/errs"
	"quicket/booking-service/pkg/util"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	FindByID(id uint) error
	Create(ctx context.Context, req *CreateBookingRequest, userPublicID string) (*BookingDTO, error)
}

type srv struct {
	repo RepositoryInterface
	eventReader EventReader
	usrReader UserReader
	logger zerolog.Logger
}

func Newsrv(repo RepositoryInterface, eventReader EventReader, usrReader UserReader, logger zerolog.Logger) *srv {
	return &srv{
		repo: repo,
		eventReader: eventReader,
		usrReader: usrReader,
		logger: logger,
	}
}

func (s *srv) FindByID(id uint) error {
	return nil
}

func (s *srv) Create(ctx context.Context, req *CreateBookingRequest, userPublicID string) (*BookingDTO, error) {
		log := s.logger.With().
		Str("event_public_id", req.EventID).
		Uint("seats", req.Seats).
		Str("user_public_id", userPublicID).
		Logger()
	
	log.Info().Msg("Create booking request")
	
	log.Debug().Msg("Checking eventID exist or not")
	ev, err := s.eventReader.GetEventDateTimeAndSeats(ctx, req.EventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewErrNotFound("event")
		}
		return nil, fmt.Errorf("booking service#create: %w", err)
	}
	now := time.Now()
	if now.After(ev.EndDate) {
		return nil, errs.NewConflictError("can not book past event")
	}

	userID, err := s.usrReader.GetUserID(ctx, userPublicID)
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

	bDTO := s.prepareBookingDTO(persisted, req.EventID, userPublicID)
	return bDTO, nil
}

func (s *srv) prepareBooking(ctx context.Context, eventID uint, userID uint, seats uint) (*Booking, error) {
	publicID, err := util.GeneratePublicID(ctx)
	if err != nil {
		return nil, fmt.Errorf("booking#create: generate public id: %w", err)
	}

	expiredAt := time.Now().Add(5 * time.Minute)

	return &Booking{
		PublicID: publicID,
		EventID: eventID,
		UserID: uint(userID),
		Seats: seats,
		ExpiredAt: expiredAt,
	}, nil
}

func (s *srv) prepareBookingDTO(booking *Booking, eventPublicID string, userPublicID string) *BookingDTO {
	return &BookingDTO{
		PublicID: booking.PublicID,
		EventPublicID: eventPublicID,
		UserID: userPublicID,
		Seats: booking.Seats,
		Status: booking.Status,
		ExpiredAt: booking.ExpiredAt,
	}
}