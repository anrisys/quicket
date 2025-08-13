package event

import (
	"context"
	"errors"
	"fmt"

	"github.com/anrisys/quicket/internal/event/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/anrisys/quicket/pkg/util"
	"github.com/rs/zerolog"
)

type EventServiceInterface interface {
	Create(ctx context.Context, req *dto.CreateEventRequest, organizer int) (*Event, error)
	FindByID(ctx context.Context, id uint) (*Event, error)
	FindByPublicID(ctx context.Context, publicID string) (*Event, error)
	eventExistsByTitle(ctx context.Context, title string) (bool, error)
	prepareEvent(ctx context.Context, req *dto.CreateEventRequest, userID int) (*Event, error)
}

type EventService struct {
	repo EventRepositoryInterface
	logger zerolog.Logger
}

func NewEventService(repo EventRepositoryInterface, logger zerolog.Logger) *EventService {
	return &EventService{
		repo: repo,
		logger: logger,
	}
}

func (s *EventService) Create(ctx context.Context, req *dto.CreateEventRequest, organizer int) (*Event, error) {
	const op = "event/service.Create"
	s.logger.Debug().Int("userId", organizer).Msg("Attempt to create event")
	s.logger.Debug().Int("userId", organizer).Msg("Checking if user can create event")
	
	// role, err := s.roleService.GetUserRole(ctx, organizer)
	// if err != nil {
	// 	return nil, fmt.Errorf("%s: failed to get user: %w", op, err)
	// }

	// if role != "organizer" && role != "admin" {
	// 	return nil, errs.ErrForbidden
	// }

	s.logger.Debug().Int("userId", organizer).Msg("Checking if webinar same title already exists")
	
	exists, err := s.eventExistsByTitle(ctx, req.Title)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if exists {
		return nil, errs.NewConflictError("event with this title already exists")
	}

	newEvent, err := s.prepareEvent(ctx, req, organizer)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	registeredEvent, err := s.repo.Create(ctx, newEvent)

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

func (s *EventService) GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*dto.EventDateTimeAndSeats, error) {
	event, err := s.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}
	return &dto.EventDateTimeAndSeats{
		ID: int(event.ID),
		AvailableSeats: event.AvailableSeats,
		EndDate: event.EndDate,
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

func (s *EventService) prepareEvent(ctx context.Context, req *dto.CreateEventRequest, userID int) (*Event, error) {
	publicID, err := util.GeneratePublicID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public ID: %w", err)
	}

	return &Event{
		PublicID: publicID,
		Title: req.Title,
		Description: &req.Description,
		OrganizerID: uint(userID),
		StartDate: req.StartDate,
		EndDate: req.EndDate,
		MaxSeats: req.MaxSeats,
		AvailableSeats: req.MaxSeats,
	}, nil
}