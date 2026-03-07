package booking

import (
	"context"
	"quicket/booking-service/internal/booking/dto"
	"quicket/booking-service/internal/domain/booking"
	"quicket/booking-service/internal/snapshot"
	"quicket/booking-service/internal/util"
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

	bookingPublicID, err := util.GeneratePublicID(ctx)
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
