package user

import "github.com/rs/zerolog"

type ServiceInterface interface {
	Register (email, password string) error
}

type service struct {
	repo *RepositoryInterface
	logger     zerolog.Logger
}

func NewService(repo *RepositoryInterface, logger *zerolog.Logger) *service {
	return &service{
		repo: repo,
		logger: *logger,
	}
}

func (s *service) Register(email, password string) error {
	return nil
}

