package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/anrisys/quicket/event-service/pkg/errs"
	redisClient "github.com/anrisys/quicket/event-service/pkg/redis"
	"github.com/anrisys/quicket/event-service/pkg/util"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type EventServiceInterface interface {
	Create(ctx context.Context, req *CreateEventRequest, userPublicID string) (*EventDTO, error)
	FindByID(ctx context.Context, id uint) (*EventDTOWithID, error) 
	FindByPublicID(ctx context.Context, publicID string) (*EventDTO, error)
	eventExistsByTitle(ctx context.Context, title string) (bool, error)
	prepareEvent(ctx context.Context, req *CreateEventRequest, userID int) (*Event, error)
}

type EventService struct {
	repo   EventRepositoryInterface
	users  UserReader
	logger zerolog.Logger
	redis *redisClient.Client
}

func NewEventService(repo EventRepositoryInterface, users UserReader, logger zerolog.Logger, redis *redisClient.Client) *EventService {
	return &EventService{
		repo:   repo,
		users:  users,
		logger: logger,
		redis: redis,
	}
}

func (s *EventService) Create(ctx context.Context, req *CreateEventRequest, userPublicID string) (*EventDTO, error) {
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

	evDTO := s.prepareEventDTO(ctx, registeredEvent)
	return evDTO, nil
}

func (s *EventService) FindByID(ctx context.Context, id uint) (*EventDTOWithID, error) {
	dbEvent, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	evDTO := s.prepareEventDTOWithID(ctx, dbEvent)
	return evDTO, nil
}

func (s *EventService) FindByPublicID(ctx context.Context, publicID string) (*EventDTO, error) {
	var event *EventDTO
	cacheKey := fmt.Sprintf("%s:publicID:%s", redisClient.EventKey, publicID)
	err := s.redis.Get(ctx, cacheKey, event)
	if err == nil {
		s.logger.Debug().Msgf("Cache hit for event with public ID: %s", publicID)
		return event, nil
	}

	if err == redis.Nil {
		s.logger.Debug().Msgf("Cache miss for event with public ID: %s", publicID)
	} else {
		s.logger.Error().Err(err).Msg("Redis Get operation failed")
	}

	dbEvent, err := s.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, err
	}

	cacheTTL := 1 * time.Hour
	err = s.redis.Set(ctx, cacheKey, dbEvent, cacheTTL)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to set event in Redis cache")
	}
	
	evDTO := s.prepareEventDTO(ctx, dbEvent)
	return evDTO, nil
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

func (s *EventService) prepareEventDTO(_ctx context.Context, ev *Event) *EventDTO {
	return &EventDTO{
		PublicID: ev.PublicID,
		Title: ev.Title,
		Description: ev.Description,
		StartDate: ev.StartDate,
		EndDate: ev.EndDate,
		MaxSeats: ev.MaxSeats,
		AvailableSeats: ev.AvailableSeats,
		CreatedAt: ev.CreatedAt,
		UpdatedAt: ev.UpdatedAt,
	}
}

func (s *EventService) prepareEventDTOWithID(_ctx context.Context, ev *Event) *EventDTOWithID {
	baseEv := s.prepareEventDTO(_ctx, ev)
	return &EventDTOWithID{
		ID: ev.ID,
		EventDTO: *baseEv,
	}
}