package event

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/anrisys/quicket/pkg/errs"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type EventRepositoryInterface interface {
	Create(ctx context.Context, event *Event) (*Event, error)
	FindByTitle(ctx context.Context, title string) (*Event, error)
	FindByID(ctx context.Context, id uint) (*Event, error)
	FindByPublicID(ctx context.Context, publicID string) (*Event, error) 
}

type EventRepository struct {
	db *gorm.DB
	logger zerolog.Logger
}

func NewEventRepository(db *gorm.DB, logger zerolog.Logger) *EventRepository {
	return &EventRepository{
		db: db,
		logger: logger,
	}
}

func (r *EventRepository) Create(ctx context.Context, event *Event) (*Event, error)  {
	err := r.db.WithContext(ctx).Create(event).Error
	if err != nil {
        if errors.Is(err, gorm.ErrDuplicatedKey) {
            return nil, errs.NewConflictError("event already exists")
        }
        if isConnectionError(err) {
            return nil, errs.NewServiceUnavailableError("database unavailable")
        }
        return nil, fmt.Errorf("failed to create event: %w", err)
	}
	return event, nil
}

func (r *EventRepository) FindByTitle(ctx context.Context, title string) (*Event, error) {
	event := &Event{}
	err := r.db.WithContext(ctx).Take(event, "title = ?", title).Error
	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errs.NewErrNotFound("event")
        }
        if isConnectionError(err) {
            return nil, errs.NewServiceUnavailableError("database unavailable")
        }
        return nil, fmt.Errorf("failed to find event by title: %w", err)
	}
	return event, nil
}

func (r *EventRepository) FindByID(ctx context.Context, id uint) (*Event, error)  {
	event := &Event{}
	err := r.db.WithContext(ctx).First(event, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewErrNotFound("event")
		}
		if isConnectionError(err) {
			return nil, errs.NewServiceUnavailableError("database unavailable")
		}
		return nil, fmt.Errorf("failed to find event by ID: %w", err)
	}
	return event, nil
}

func (r *EventRepository) FindByPublicID(ctx context.Context, publicID string) (*Event, error) {
	event := &Event{}
	err := r.db.WithContext(ctx).Take(event, "public_id = ?", publicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NewErrNotFound("event")
		}
		if isConnectionError(err) {
			return nil, errs.NewServiceUnavailableError("database unavailable")
		}
		return nil, fmt.Errorf("failed to find event by ID: %w", err)
	}
	return event, nil
}

func isConnectionError(err error) bool {
    // Implement proper connection error detection
    return strings.Contains(err.Error(), "connection refused") || 
           errors.Is(err, context.DeadlineExceeded)
}