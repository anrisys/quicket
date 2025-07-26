package user

import (
	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *User) error
	FindById(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

type UserRepository struct {
	db 		*gorm.DB
	logger 	zerolog.Logger
}

func NewUserRepository(db *gorm.DB, logger zerolog.Logger) *UserRepository {
	return &UserRepository{
		db: db,
		logger: logger,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return nil
	}
	return nil
}

func (r *UserRepository) FindById(ctx context.Context, id int) (*User, error)  {
	// should check if there is error or not? 
	return nil, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error)  {
	// implementation
	return nil, nil
}

