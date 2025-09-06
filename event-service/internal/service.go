package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/anrisys/quicket/event-service/pkg/errs"
	"github.com/anrisys/quicket/event-service/pkg/util"
	"github.com/rs/zerolog"
)

type EventServiceInterface interface {
	Create(ctx context.Context, req *CreateEventRequest, userPublicID string) (*Event, error)
	FindByID(ctx context.Context, id uint) (*Event, error)
	FindByPublicID(ctx context.Context, publicID string) (*Event, error)
	eventExistsByTitle(ctx context.Context, title string) (bool, error)
	prepareEvent(ctx context.Context, req *CreateEventRequest, userID int) (*Event, error)
}

type EventService struct {
	repo   EventRepositoryInterface
	users  UserReader
	logger zerolog.Logger
}

func NewEventService(repo EventRepositoryInterface, users UserReader, logger zerolog.Logger) *EventService {
	return &EventService{
		repo:   repo,
		users:  users,
		logger: logger,
	}
}

func (s *EventService) Create(ctx context.Context, req *CreateEventRequest, userPublicID string) (*Event, error) {
	log := s.logger.With().
		Str("user_id", userPublicID).
		Logger()

	log.Info().Msg("Create new event")
	usr, err := s.users.FindUserByPublicID(ctx, userPublicID)
	if err != nil {
		return nil, fmt.Errorf("event/service#create: %w", err)
	}

	exists, err := s.eventExistsByTitle(ctx, req.Title)
	if err != nil {
		return nil, fmt.Errorf("event/service#create: %w", err)
	}

	if exists {
		return nil, errs.NewConflictError("event with this title already exists")
	}

	newEv, err := s.prepareEvent(ctx, req, usr.ID)
	if err != nil {
		return nil, err
	}

	registeredEvent, err := s.repo.Create(ctx, newEv)

	if err != nil {
		return nil, fmt.Errorf("event service#create: %w", err)
	}

	return registeredEvent, err
}

func (s *EventService) FindByID(ctx context.Context, id uint) (*Event, error) {
	event, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventService) FindByPublicID(ctx context.Context, publicID string) (*Event, error) {
	event, err := s.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *EventService) GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*EventDateTimeAndSeats, error) {
	event, err := s.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}
	return &EventDateTimeAndSeats{
		ID:             int(event.ID),
		AvailableSeats: event.AvailableSeats,
		EndDate:        event.EndDate,
	}, nil
}

func (s *EventService) eventExistsByTitle(ctx context.Context, title string) (bool, error) {
	_, err := s.repo.FindByTitle(ctx, title)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, errs.ErrNotFound) {
		return false, nil
	}
	return false, err
}

func (s *EventService) prepareEvent(ctx context.Context, req *CreateEventRequest, userID int) (*Event, error) {
	publicID, err := util.GeneratePublicID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public ID: %w", err)
	}

	return &Event{
		PublicID:       publicID,
		Title:          req.Title,
		Description:    &req.Description,
		OrganizerID:    uint(userID),
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		MaxSeats:       req.MaxSeats,
		AvailableSeats: req.MaxSeats,
	}, nil
}