package eventsnapshot

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Repository interface {
	Find(ctx context.Context, id uint) (*EventSnapshot, error)
	Create(ctx context.Context, ev *EventSnapshot) error
	Update(ctx context.Context, ev *EventSnapshot) error
	UpdateSeats(ctx context.Context, evID uint, available_seats, version int) error
	Delete(ctx context.Context, id uint) error
	GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*EventDateTimeAndSeats, error)
}

type EvSnapshotRepo struct {
	db *gorm.DB
	logger zerolog.Logger
}

func NewEvSnapshotRepo(db *gorm.DB, logger zerolog.Logger) *EvSnapshotRepo {
	return &EvSnapshotRepo{
		db: db,
		logger: logger,
	}
}

func (r *EvSnapshotRepo) Find(ctx context.Context, id uint) (*EventSnapshot, error) {
	var ev EventSnapshot
	err := r.db.First(&ev, id).Error
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

func (r *EvSnapshotRepo) Create(ctx context.Context, ev *EventSnapshot) error {
	err := r.db.WithContext(ctx).Create(ev).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *EvSnapshotRepo) Update(ctx context.Context, ev *EventSnapshot) error {
	err := r.db.WithContext(ctx).Save(ev).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *EvSnapshotRepo) UpdateSeats(ctx context.Context, evID uint, available_seats, version int) error {
	var ev EventSnapshot
	err := r.db.First(&ev, evID).Error
	if err != nil {
		return fmt.Errorf("failed to select event: %w", err)
	}
	err = r.db.Model(&ev).Select("available_seats", "version").Updates(EventSnapshot{AvailableSeats: uint64(available_seats), Version: ev.Version + 1}).Error
	if err != nil {
		return fmt.Errorf("failed to update an event snapshot: %w", err)
	}
	return nil
}

func (r *EvSnapshotRepo) Delete(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&EventSnapshot{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *EvSnapshotRepo) GetEventDateTimeAndSeats(ctx context.Context, publicID string) (*EventDateTimeAndSeats, error) {
	var ev *EventDateTimeAndSeats
	err := r.db.WithContext(ctx).Select("id", "available_seats", "end_date").Take(ev, "public_id = ?", publicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		r.logger.Error().Err(err).
			Str("event:public_id", publicID).
			Msg("event not found")
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}
	return ev, nil
}