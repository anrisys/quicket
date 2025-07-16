package user

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type RepositoryInterface interface {
	Create(user *User) error
	FindById(id int) (*User, error)
	FindByEmail(email string) (*User, error)
}

type repository struct {
	db *gorm.DB
	logger zerolog.Logger
}

func NewRepository(db *gorm.DB, logger zerolog.Logger) *repository {
	return &repository{
		db: db,
		logger: logger,
	}
}

func (r *repository) Create(user *User) error {
	// implementation
	return nil
}

func (r *repository) FindById(id int) (*User, error)  {
	// implementation
	return nil, nil
}

func (r *repository) FindByEmail(id int) (*User, error)  {
	// implementation
	return nil, nil
}

