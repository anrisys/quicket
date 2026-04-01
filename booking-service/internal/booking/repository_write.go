package booking

import (
	"context"
	"fmt"
	"quicket/booking-service/internal/domain/booking"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

/*
BEST PRACTICE:
1. LOGGING: wrap error found in the repo layer and let it bubble up. don't log the error in repo layer since it causes log noises

2. ERROR WRAPPING CONVENTION:
"<layer>.<method>: <what you were doing>: %w"
e.g. "userRepo.GetByID: scan result: %w"
*/
type BookingWriteRepository interface {
	Create(ctx context.Context, b *Booking) error
	ExpirePending(ctx context.Context, limit uint64) error
	Cancel(ctx context.Context, bookingID uint64) error
	ConfirmSuccess(ctx context.Context, bookingID uint64) error
	Fail(ctx context.Context, bookingID uint64) error
}

type BookingWriteRepoImpl struct {
	db *gorm.DB
	l  zerolog.Logger
}

func NewBookingWriteRepoImpl(db *gorm.DB, l zerolog.Logger) *BookingWriteRepoImpl {
	return &BookingWriteRepoImpl{db: db, l: l}
}

func (r *BookingWriteRepoImpl) Create(ctx context.Context, b *Booking) error {
	return r.db.WithContext(ctx).Create(b).Error
}

// Target SQL: UPDATE bookings SET status = 'expired' WHERE status = 'pending' and expired_at <= NOW() ORDER BY id desc LIMIT 100
func (r *BookingWriteRepoImpl) ExpirePending(ctx context.Context, limit uint64) error {
	err := r.db.WithContext(ctx).
		Model(&Booking{}).
		Where("status = ?", booking.BookingStatusPending).
		Where("expired_at <= ?", time.Now().UTC()).
		Order("id desc").
		Limit(int(limit)).
		Update("status", booking.BookingStatusExpired).
		Error

	if err != nil {
		return fmt.Errorf("repository_write.ExpirePending: query failed: %w", err)
	}
	return nil
}

func (r *BookingWriteRepoImpl) Cancel(ctx context.Context, bookingID uint64) error {
	err := r.db.WithContext(ctx).
		Model(&Booking{}).
		Where("id = ?", bookingID).
		Where("status = ?", booking.BookingStatusPending).
		Update("status", booking.BookingStatusCancelled).
		Error

	if err != nil {
		return fmt.Errorf("repository_write.Cancel: query failed: %w", err)
	}
	return nil
}

func (r *BookingWriteRepoImpl) ConfirmSuccess(ctx context.Context, bookingID uint64) error {
	err := r.db.WithContext(ctx).
		Model(&Booking{}).
		Where("id = ?", bookingID).
		Where("status = ?", booking.BookingStatusSuccess).
		Update("status", booking.BookingStatusPending).
		Error

	if err != nil {
		return fmt.Errorf("repository_write.Confirm: query failed: %w", err)
	}
	return nil
}

func (r *BookingWriteRepoImpl) Fail(ctx context.Context, bookingID uint64) error {
	err := r.db.WithContext(ctx).
		Model(&Booking{}).
		Where("id = ?", bookingID).
		Where("status = ?", booking.BookingStatusFailed).
		Update("status", booking.BookingStatusPending).
		Error

	if err != nil {
		return fmt.Errorf("repository_write.Fail: query failed: %w", err)
	}
	return nil
}
