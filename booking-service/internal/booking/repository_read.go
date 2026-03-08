package booking

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BookingReadRepo interface {
	ListAdminBookings() ([]BookingRow, error)
	RetrieveExpiredBookings(ctx context.Context, limit uint64) ([]ExpiredBookingRow, error)
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

func (r *BookingReadRepoImpl) RetrieveExpiredBookings(
	ctx context.Context,
	limit uint64,
) ([]ExpiredBookingRow, error) {

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
		return nil, fmt.Errorf("repository_read.RetrieveExpiredBookings: query failed: %w", err)
	}

	return rows, nil
}
