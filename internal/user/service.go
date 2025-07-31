package user

import (
	"context"
	"fmt"

	"github.com/anrisys/quicket/internal/user/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/anrisys/quicket/pkg/security"
	"github.com/rs/zerolog"
)

type UserServiceInterface interface {
	Register (ctx context.Context, req *dto.RegisterUserRequest) error
	Login (ctx context.Context, req *dto.LoginUserRequest) (*dto.UserDTO, error)
}

type UserService struct {
	repo 			UserRepositoryInterface
	logger  		zerolog.Logger
	accountSecurity security.AccountSecurityInterface
}

func NewUserService(
	repo UserRepositoryInterface, 
	logger zerolog.Logger, 
	accountSecurity security.AccountSecurityInterface,
) *UserService {
	return &UserService{
		repo: repo,
		logger: logger,
		accountSecurity: accountSecurity,
	}
}

func (s *UserService) Register(ctx context.Context, req *dto.RegisterUserRequest) error {
	s.logger.Debug().Str("email", req.Email).Msg("Attempt to register a user")
	
	s.logger.Debug().Msgf("Checking if emai exists %s", req.Email)
	emailRegistered := s.repo.EmailExists(ctx, req.Email)

	if emailRegistered {
		return errs.NewConflictError("email already registered")
	}
	
	hashedPassword, err := s.accountSecurity.HashPassword(req.Password)
	if err != nil {
		return errs.NewInternalError("failed to hash password")
	}

	publicID, err := s.accountSecurity.GeneratePublicID()
	if err != nil {
		return errs.NewInternalError("failed to generate user public id")
	}

	user := &User{
		PublicID: publicID,
		Email: req.Email,
		Password: hashedPassword,
		Role: "user",
	}
	
	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("user service: create operation: %w", err)
	}

	s.logger.Info().Msgf("New user registered %s", req.Email)
	return nil
}

func (s *UserService) Login(ctx context.Context, req *dto.LoginUserRequest) (*dto.UserDTO, error) {
	s.logger.Debug().Ctx(ctx).Str("email", req.Email).Msg("Attempt to login")

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("user service: login: %w", err)
	}

	passwordMatch := s.accountSecurity.CheckPasswordHash(req.Password, user.Password)
	if !passwordMatch {
		return nil, fmt.Errorf("email or password is wrong %w",
			errs.NewAppError(400, "INVALID_DATA", "email or password is wrong"),
		)
	}

	userDto := dto.UserDTO{
		Email: user.Email,
		PublicID: user.PublicID,
		Role: user.Role,
	}

	s.logger.Info().Ctx(ctx).Str("userId", userDto.PublicID).Msg("User login")
	return &userDto, nil
}
