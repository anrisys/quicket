package user

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(user *User) error
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
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

func (r *UserRepository) Create(user *User) error {
	// implementation
	return nil
}

func (r *UserRepository) FindById(id int) (*User, error)  {
	// implementation
	return nil, nil
}

func (r *UserRepository) FindByEmail(email string) (*User, error)  {
	// implementation
	return nil, nil
}

