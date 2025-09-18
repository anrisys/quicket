package eventsnapshot

import (
	"context"

	"github.com/rs/zerolog"
)

type Service interface {
	CreateSnapshot(ctx context.Context, ev *EventSnapshot) error
	UpdateSnapshot(ctx context.Context, ev *EventSnapshot) error
	UpdateSeatsSnapshot(ctx context.Context, evID uint, available_seats, version int) error
	DeleteSnapshot(ctx context.Context, id uint) error
	GetEventSnapshotDateTimeAndSeats(ctx context.Context, publicID string) (*EventDateTimeAndSeats, error)
}

type srv struct {
	repo Repository
	logger zerolog.Logger
}

func NewEvSnapshotSrv(repo Repository, logger zerolog.Logger) *srv {
	return &srv{
		repo: repo,
		logger: logger,
	}
}

func (s *srv) CreateSnapshot(ctx context.Context, ev *EventSnapshot) error {
	err := s.repo.Create(ctx, ev)
	return err
}

func (s *srv) UpdateSnapshot(ctx context.Context, ev *EventSnapshot) error {
	err := s.repo.Update(ctx, ev)
	return err
}

func (s *srv) UpdateSeatsSnapshot(ctx context.Context, evID uint, available_seats, version int) error {
	err := s.repo.UpdateSeats(ctx, evID, available_seats, version)
	return err
}

func (s *srv) DeleteSnapshot(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	return err
}

func (s *srv) GetEventSnapshotDateTimeAndSeats(ctx context.Context, publicID string) (*EventDateTimeAndSeats, error) {
	ev, err := s.repo.GetEventDateTimeAndSeats(ctx, publicID)
	if err != nil {
		return nil, err
	}
	return ev, nil
}