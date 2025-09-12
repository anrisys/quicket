package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type EventReader interface {
	GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*EventDateTimeAndSeats, error)
}

type EvReader struct {
	db *gorm.DB
	logger zerolog.Logger
}

func NewEvReader(db *gorm.DB, logger zerolog.Logger) *EvReader {
	return &EvReader{
		db: db,
		logger: logger,
	}
}

func (e *EvReader) GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*EventDateTimeAndSeats, error) {
	var ev *EventDateTimeAndSeats
	err := e.db.WithContext(ctx).Select("id", "available_seats", "end_date").Take(ev, "public_id = ?", publicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		e.logger.Error().Err(err).
			Str("event:public_id", publicID).
			Msg("event not found")
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}
	return ev, nil
}