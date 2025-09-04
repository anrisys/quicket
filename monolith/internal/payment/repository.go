package payment

import (
	"context"
	"errors"
	"fmt"

	commonDTO "github.com/anrisys/quicket/internal/dto"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BookingRow struct {
	ID uint
	PublicID string
	Status string
}

type Repository interface {
	CreatePaymentAndUpdateBookingStatus(ctx context.Context, p *Payment) (*commonDTO.PaymentDTO, error)
}

type GormRepository struct {
	db *gorm.DB
	logger zerolog.Logger
}

func NewRepository(db *gorm.DB, logger zerolog.Logger) *GormRepository {
	return &GormRepository{
		db: db,
		logger: logger,
	}
}

func (r *GormRepository) CreatePaymentAndUpdateBookingStatus(ctx context.Context, p *Payment) (*commonDTO.PaymentDTO, error) {
	var b BookingRow
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("bookings").
			Where("id = ?", p.BookingID).
			Select("id", "public_id", "status").
			Take(&b).Error; err != nil {
			
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return ErrBookingNotFound
				}
				
				r.logger.Error().Err(err).
					Uint("booking_id", p.BookingID).
					Msg("select booking failed")
				return fmt.Errorf("%w: %v", ErrDB, err)
		}

		if err := tx.Table("payments").Create(p).Error; err != nil {
			r.logger.Error().Err(err).
				Uint("booking_id", p.BookingID).
				Uint("user_id", p.UserID).
				Msg("insert payment failed")
			return fmt.Errorf("%w: %v", ErrDB, err)
		}

		
		if err := tx.Table("bookings").Where("id = ?", b.ID).Update("status", p.Status).Error; err != nil {
			r.logger.Error().Err(err).
				Uint("booking_id", b.ID).
				Msg("failed to update booking status")
			return fmt.Errorf("%w: %v", ErrDB, err)
		}

		return nil
	})

	return &commonDTO.PaymentDTO{
		PublicID: p.PublicID,
		Amount: p.Amount,
		Status: p.Status,
		BookingID: b.PublicID,
	}, err
}