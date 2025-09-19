package usersnapshot

import (
	"context"

	"github.com/rs/zerolog"
)

type Service interface {
	GetUserSnapshotID(ctx context.Context, publicID string) (*uint, error)
	CreateUserSnapshot(ctx context.Context, userID uint, userPublicID string) error
	DeleteUserSnapshot(ctx context.Context, id uint) error
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

func (s *srv) CreateUserSnapshot(ctx context.Context, userID uint, userPublicID string) error {
	usr := UserSnapshot{
		ID: userID,
		PublicID: userPublicID,
	}
	return s.repo.Create(ctx, &usr)
}

func (s *srv) DeleteUserSnapshot(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}