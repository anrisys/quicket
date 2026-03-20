package booking

import (
	"context"
	"errors"
	"fmt"
	"quicket/booking-service/internal/domain/booking"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BookingReadRepo interface {
	ListAdminBookings() ([]BookingRow, error)
	FindExpiredPending(ctx context.Context, limit uint64) ([]ExpiredBookingRow, error)
	FindForCancellation(ctx context.Context, bookingID string) (*BookingCancellationRow, error)
}

type BookingReadRepoImpl struct {
	db     *gorm.DB
	logger zerolog.Logger
}

func NewBookingReadRepo(db *gorm.DB, logger zerolog.Logger) *BookingReadRepoImpl {
	return &BookingReadRepoImpl{db: db, logger: logger}
}

func (r *BookingReadRepoImpl) ListAdminBookings() ([]BookingRow, error) {
	return nil, nil
}

func (r *BookingReadRepoImpl) FindExpiredPending(ctx context.Context, limit uint64) ([]ExpiredBookingRow, error) {

	var rows []ExpiredBookingRow

	err := r.db.WithContext(ctx).
		Table("bookings").
		Select("bookings.id, event_snapshot.public_id as event_public_id, bookings.seats").
		Joins("left join event_snapshot on bookings.event_id = event_snapshot.id").
		Where("bookings.status = ?", "pending").
		Where("bookings.expired_at <= ?", time.Now().UTC()).
		Order("bookings.id desc").
		Limit(int(limit)).
		Scan(&rows).
		Error

	if err != nil {
		return nil, fmt.Errorf("repository_read.FindExpiredBooking: query failed: %w", err)
	}

	return rows, nil
}

func (r *BookingReadRepoImpl) FindForCancellation(ctx context.Context, bookingID string) (*BookingCancellationRow, error) {
	var book BookingCancellationRow

	err := r.db.WithContext(ctx).
		Table("bookings").
		Select("bookings.id, bookings.public_id, bookings.status, bookings.seats, bookings.event_id as event_public_id, bookings.user_id as user_public_id").
		Where("bookings.public_id = ?", bookingID).
		Limit(1).
		Scan(&book).
		Error

	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, booking.ErrBookingNotFound
		}
		return nil, fmt.Errorf("repository_read.FindForCancellation: query failed: %w", err)
	}

	return &book, nil
}
