package user

import "github.com/rs/zerolog"

type UserServiceInterface interface {
	Register (email, password string) error
}

type UserService struct {
	repo 	UserRepositoryInterface
	logger  zerolog.Logger
}

func NewUserService(repo UserRepositoryInterface, logger zerolog.Logger) *UserService {
	return &UserService{
		repo: repo,
		logger: logger,
	}
}

func (s *UserService) Register(email, password string) error {
	return nil
}

