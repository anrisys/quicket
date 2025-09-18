package usersnapshot

import (
	"context"
	"errors"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserID(ctx context.Context, publicID string) (*uint, error)
}

type repo struct {
	db *gorm.DB
	logger zerolog.Logger
}

func NewRepo(db *gorm.DB, logger zerolog.Logger) *repo {
	return &repo{
		db: db,
		logger: logger,
	}
}

func (r *repo) GetUserID(ctx context.Context, publicID string) (*uint, error) {
	var usr *uint
	err := r.db.WithContext(ctx).Select("id").Take(usr, "public_id = ?", publicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		r.logger.Error().Err(err).
			Str("user:public_id:", publicID).
			Msg("user not found")
		return nil, fmt.Errorf("%w: %v", ErrDB, err)
	}
	return usr, nil
}