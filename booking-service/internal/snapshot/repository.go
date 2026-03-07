package snapshot

import (
	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

/*
|--------------------------------------------------------------------------
| Event Snapshot Repository
|--------------------------------------------------------------------------
*/
type EventIDs struct {
	ID       uint64
	PublicID string
}

type EventSnapshotRepository interface {
	Create(ctx context.Context, ev EventSnapshot) error
	FindIDsByPublicID(ctx context.Context, publicID string) (*EventIDs, error)
	FindEventSeatPrice(ctx context.Context, id uint64) (float64, error)
}

type EvSnapshotRepoImpl struct {
	db *gorm.DB
	l  zerolog.Logger
}

func NewEvSnapshotRepoImpl(db *gorm.DB, l zerolog.Logger) *EvSnapshotRepoImpl {
	return &EvSnapshotRepoImpl{db: db, l: l}
}

func (r *EvSnapshotRepoImpl) Create(ctx context.Context, ev EventSnapshot) error {
	err := r.db.Create(&ev).Error
	r.l.Error().Err(err)
	return err
}

func (r *EvSnapshotRepoImpl) FindIDsByPublicID(ctx context.Context, publicID string) (*EventIDs, error) {
	var evIDs EventIDs
	err := r.db.Model(&EventSnapshot{}).
		Where("public_id = ?", publicID).
		Select(&evIDs).Error
	if err != nil {
		r.l.Error().Err(err)
		return nil, err
	}
	return &evIDs, nil
}

func (r *EvSnapshotRepoImpl) FindEventSeatPrice(ctx context.Context, id uint64) (float64, error) {
	var price float64
	err := r.db.Model(&EventSnapshot{}).
		Where("id = ?", id).
		Select(price).Error
	if err != nil {
		r.l.Error().Err(err)
		return 0, err
	}
	return price, nil
}

/*
|--------------------------------------------------------------------------
| User Snapshot Repository
|--------------------------------------------------------------------------
*/
type UserSnapshotRepository interface {
	Create(ctx context.Context, us UserSnapshot) error
	FindUserPrimaryIDByPublicID(ctx context.Context, publicID string) (uint64, error)
}

type UserSnapshotRepoImpl struct {
	db *gorm.DB
	l  zerolog.Logger
}

func NewUserSnapshotRepoImpl(db *gorm.DB, l zerolog.Logger) *UserSnapshotRepoImpl {
	return &UserSnapshotRepoImpl{db: db, l: l}
}

func (r *UserSnapshotRepoImpl) Create(ctx context.Context, us UserSnapshot) error {
	err := r.db.Create(&us).Error
	r.l.Error().Err(err)
	return err
}

func (r *UserSnapshotRepoImpl) FindUserPrimaryIDByPublicID(ctx context.Context, publicID string) (uint64, error) {
	var userPrimaryID uint64
	err := r.db.Model(&UserSnapshot{}).
		Where("public_id = ?", publicID).
		Select("id").
		First(userPrimaryID).
		Error
	r.l.Error().Err(err)
	if err != nil {
		return 0, err
	}
	return userPrimaryID, nil
}
