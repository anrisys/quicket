package booking

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(ctx context.Context, b *Booking) (*Booking, error)
}

type eventRow struct {
	ID uint
	AvailableSeats uint64
}

type GormRepository struct {
	db *gorm.DB
	logger zerolog.Logger
}

func NewGormRepository(db *gorm.DB, logger zerolog.Logger) *GormRepository {
	return &GormRepository{
		db: db,
		logger: logger,
	}
}

func (r *GormRepository) Create(ctx context.Context, b *Booking) (*Booking, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ev eventRow
		if err := tx.Table("events").
			Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).
			Select("id", "available_seats").
			Where("id = ?", b.EventID).
			Take(&ev).Error; err != nil {
			
				if err == gorm.ErrRecordNotFound {
					return ErrEventNotFound
				}

				r.logger.Error().Err(err).
					Uint("event_id", b.EventID).
					Msg("lock/select event failed")
				return fmt.Errorf("%w: %v", ErrDB, err)
		}

		if ev.AvailableSeats < uint64(b.Seats) {
			return ErrNotEnoughSeats
		}

		if err := tx.Create(b).Error; err != nil  {
			r.logger.Error().Err(err).
				Uint("event_id", b.EventID).
				Uint("user_id", b.UserID).
				Msg("insert booking failed")
			return fmt.Errorf("%w: %v", ErrDB, err)
		}

		if err := tx.Table("events").
			Where("id = ?", b.EventID).
			Update("available_seats", gorm.Expr("available_seats - ?", b.Seats)).
			Error; err != nil {
			r.logger.Error().Err(err).
				Uint("event_id", b.EventID).
				Uint("seats", b.Seats).
				Msg("deduct seats failed")
			return fmt.Errorf("%w: %v", ErrDB, err)
		}

		return nil
	})

	return b, err
}