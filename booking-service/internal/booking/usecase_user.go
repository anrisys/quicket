package booking

import (
	"context"
	"fmt"
	"quicket/booking-service/internal/booking/dto"
	"quicket/booking-service/internal/domain/booking"
	"quicket/booking-service/internal/helper"
	"quicket/booking-service/internal/snapshot"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type UserUsecase struct {
	wr  BookingWriteRepository
	rr  BookingReadRepo
	l   zerolog.Logger
	ecs EventClientService
	ep  EventPublisher
	ess *snapshot.EventSnapshotService
	uss *snapshot.UserSnapshotService
}

func NewUserUsecase(
	wr BookingWriteRepository,
	rr BookingReadRepo,
	l zerolog.Logger,
	ecs EventClientService,
	ep EventPublisher,
	ess *snapshot.EventSnapshotService,
	uss *snapshot.UserSnapshotService,
) *UserUsecase {
	return &UserUsecase{wr: wr, rr: rr, l: l, ecs: ecs, ep: ep, ess: ess, uss: uss}
}

func (uu *UserUsecase) CreateBooking(ctx context.Context, data *dto.CreateBookingRequest, user string) (*dto.CreateBookingData, error) {
	// Check if event exists
	evIDs, err := uu.ess.FindIDsByPublicID(ctx, data.EventPublicID)
	if err != nil {
		return nil, err
	}
	if evIDs == nil {
		return nil, NewAppError(400, "EVENT_NOT_EXISTS", "Event not exist", booking.ErrEventNotFound)
	}

	userPrimaryID, err := uu.uss.GetUserPrimaryID(ctx, user)
	if err != nil {
		return nil, err
	}

	// Check available seats to event-service
	price, err := uu.ecs.ReserveSeats(ctx, data.EventPublicID, data.Seats)
	if err != nil {
		return nil, err
	}

	bookingPublicID, err := helper.GeneratePublicID(ctx)
	if err != nil {
		return nil, err
	}

	totalPrice := price * float64(data.Seats)

	expiredAt := time.Now().Add(10 * time.Minute)

	// Create booking entity
	bookingEntity := Booking{
		PublicID:      bookingPublicID,
		EventID:       evIDs.ID,
		UserID:        userPrimaryID,
		Seats:         data.Seats,
		TotalPrice:    totalPrice,
		Status:        booking.BookingStatusPending,
		PaymentMethod: booking.PaymentMethod(data.PaymentMethod),
		Channel:       booking.Channel(data.Channel),
		ExpiredAt:     expiredAt,
		Metadata:      nil,
	}

	// return newly created booking
	err = uu.wr.Create(ctx, &bookingEntity)
	if err != nil {
		return nil, err
	}

	result := dto.CreateBookingData{
		BookingPublicID: bookingEntity.PublicID,
		Status:          string(bookingEntity.Status),
		ExpiredAt:       expiredAt.String(),
		TotalPrice:      strconv.FormatFloat(totalPrice, 'f', 2, 64),
		Currency:        bookingEntity.Currency,
	}

	return &result, nil
}

func (uu *UserUsecase) CancelBooking(ctx context.Context, data *dto.CancelBookingRequest, user string) (*dto.CancelledBookingData, error) {
	b, err := uu.rr.FindForCancellation(ctx, data.BookingID)
	if err != nil {
		return nil, fmt.Errorf("usecase_user.CancelBooking: %w", err)
	}

	// Validate the user request and booking
	if b.UserPublicID != user {
		return nil, NewAppError(403, CodeForbidden, "You don't have permission for this action", nil)
	}
	if b.BookingStatus == booking.BookingStatusCancelled {
		return nil, NewAppError(409, CodeForbidden, "Booking is already cancelled", nil)
	}
	if b.BookingStatus == booking.BookingStatusExpired {
		return nil, NewAppError(409, CodeForbidden, "Booking is already expired", nil)
	}
	if b.BookingStatus == booking.BookingStatusFailed {
		return nil, NewAppError(409, CodeForbidden, "Booking is already failed", nil)
	}
	if b.BookingStatus == booking.BookingStatusSuccess {
		return nil, NewAppError(409, CodeForbidden, "Booking is already success", nil)
	}

	// Cancel the booking
	err = uu.wr.Cancel(ctx, b.ID)
	if err != nil {
		return nil, fmt.Errorf("usecase_user.CancelBooking: cancel failed: %w", err)
	}

	// Release seats
	uu.ep.ReleaseEventSeats(ctx, &BookingReleaseSeats{EventPublicID: b.EventPublicID, Seats: b.BookingSeats})

	return &dto.CancelledBookingData{BookingPublicID: b.BookingPublicID, Status: string(booking.BookingStatusCancelled)}, nil
}
