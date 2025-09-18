package usersnapshot

import (
	"context"

	"github.com/rs/zerolog"
)

type Service interface {
	GetUserSnapshotID(ctx context.Context, publicID string) (*uint, error)
}

type srv struct {
	repo Repository
	logger zerolog.Logger
}

func NewSrv(repo Repository, logger zerolog.Logger) *srv {
	return &srv{
		repo: repo,
		logger: logger,
	}
}

func (s *srv) GetUserSnapshotID(ctx context.Context, publicID string) (*uint, error) {
	id, err := s.repo.GetUserID(ctx, publicID)
	if err != nil {
		return nil, err
	}
	return id, nil
}