package booking

import (
	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type BookingWriteRepository interface {
	Create(ctx context.Context, b *Booking) error
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
