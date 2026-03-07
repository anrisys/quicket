package booking

import "github.com/rs/zerolog"

type AdminUsecase struct {
	brr BookingReadRepo
	l   zerolog.Logger
}

func NewAdminUsecase(brr BookingReadRepo, l zerolog.Logger) *AdminUsecase {
	return &AdminUsecase{brr: brr, l: l}
}
