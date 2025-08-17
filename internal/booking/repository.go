package booking

import (
	"context"
	"errors"
	"fmt"

	commonDTO "github.com/anrisys/quicket/internal/dto"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Create(ctx context.Context, b *Booking) (*Booking, error)
	FindSimpleDTO(ctx context.Context, publicID string) (*commonDTO.SimpleBookingDTO, error)
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

func (r *GormRepository) FindSimpleDTO(ctx context.Context, publicID string) (*commonDTO.SimpleBookingDTO, error) {
	var dto commonDTO.SimpleBookingDTO
	if err := r.db.WithContext(ctx).Select("id, user_id, total_price AS amount").Where("public_id = ?", publicID).Take(dto).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}
		r.logger.Error().Err(err).
			Str("booking_public_id", publicID).
			Msg("find booking failed")
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}
	return &dto, nil
}

func (r *GormRepository) Create(ctx context.Context, b *Booking) (*Booking, error) {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ev eventRow
		if err := tx.Table("events").
			Clauses(clause.Locking{Strength: "UPDATE", Options: "NOWAIT"}).
			Select("id, available_seats").
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