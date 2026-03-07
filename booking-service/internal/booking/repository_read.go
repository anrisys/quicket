package booking

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BookingReadRepo interface {
	ListAdminBookings() ([]BookingRow, error)
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
