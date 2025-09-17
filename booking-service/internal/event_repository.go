package internal

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type EventSnapshotRepository interface {
	CreateSnapshot(ctx context.Context, ev *EventSnapshot) error
	UpdateSnapshot(ctx context.Context, ev *EventSnapshot) error
	UpdateSeats(ctx context.Context, evID uint, available_seats, version int) error
	DeleteSnapshot(ctx context.Context, id uint) error
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

func (r *EvSnapshotRepo) CreateSnapshot(ctx context.Context, ev *EventSnapshot) error {
	err := r.db.WithContext(ctx).Create(ev).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *EvSnapshotRepo) UpdateSnapshot(ctx context.Context, ev *EventSnapshot) error {
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

func (r *EvSnapshotRepo) DeleteSnapshot(ctx context.Context, id uint) error {
	err := r.db.WithContext(ctx).Delete(&EventSnapshot{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

