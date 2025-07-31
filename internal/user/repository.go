package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/anrisys/quicket/pkg/errs"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *User) error
	FindById(ctx context.Context, id int) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	EmailExists(ctx context.Context, email string) bool
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
        r.logger.Error().Err(err).Msg("DB operation failed")
        
        switch {
        case errors.Is(err, gorm.ErrDuplicatedKey):
            return errs.NewConflictError("user already exists")
		case isConnectionError(err):
            return errs.NewServiceUnavailableError("database unavailable")
        default:
            return errs.NewAppError(500, "DB_OPERATION_FAILED", "Database operation failed", err)
        }
    }
	return nil
}

func (r *UserRepository) FindById(ctx context.Context, id int) (*User, error)  {
	// should check if there is error or not? 
	return nil, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*User, error)  {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		r.logger.Error().Ctx(ctx).Msg("DB operation failed")

		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			return nil, fmt.Errorf("%w", errs.ErrNotFound)
		case isConnectionError(err):
			return nil, fmt.Errorf("db connection: %w", errs.NewServiceUnavailableError("database unavailable"))
		default: 
			return nil, fmt.Errorf("db query: %w", errs.NewAppError(500, "DB_OPERATION_FAILED", "Database operation failed", err))
		}
	}
	return &user, nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) bool {
	var count int64

	r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count)
	
	return count > 0
}

func isConnectionError(err error) bool {
    // Implement proper connection error detection
    return strings.Contains(err.Error(), "connection refused") || 
           errors.Is(err, context.DeadlineExceeded)
}
