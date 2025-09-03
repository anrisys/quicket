package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/anrisys/quicket/user-service/pkg/errs"
	"github.com/anrisys/quicket/user-service/pkg/security"
	"github.com/anrisys/quicket/user-service/pkg/token"
	"github.com/rs/zerolog"
)

type UserServiceInterface interface {
	Register(ctx context.Context, req *RegisterUserRequest) error
	Login(ctx context.Context, req *LoginUserRequest) (*LoginUserDTO, error)
	FindUserById(ctx context.Context, id int) (*UserDTO, error)
	FindUserByPublicID(ctx context.Context, publicID string) (*UserDTO, error)
	GetUserPrimaryID(ctx context.Context, publicID string) (*uint, error)
}

type UserService struct {
	repo            UserRepositoryInterface
	logger          zerolog.Logger
	accountSecurity security.AccountSecurityInterface
	tokenGenerator  token.TokenGeneratorInterface
}

func NewUserService(
	repo UserRepositoryInterface,
	logger zerolog.Logger,
	accountSecurity security.AccountSecurityInterface,
	tokenGenerator token.TokenGeneratorInterface,
) *UserService {
	return &UserService{
		repo:            repo,
		logger:          logger,
		accountSecurity: accountSecurity,
		tokenGenerator:  tokenGenerator,
	}
}

func (s *UserService) Register(ctx context.Context, req *RegisterUserRequest) error {
	s.logger.Debug().Str("email", req.Email).Msg("Attempt to register a user")

	s.logger.Debug().Msgf("Checking if emai exists %s", req.Email)
	emailRegistered := s.repo.EmailExists(ctx, req.Email)

	if emailRegistered {
		return errs.NewConflictError("email already registered")
	}

	hashedPassword, err := s.accountSecurity.HashPassword(ctx, req.Password)
	if err != nil {
		return errs.NewInternalError("failed to hash password")
	}

	publicID, err := s.accountSecurity.GeneratePublicID(ctx)
	if err != nil {
		return errs.NewInternalError("failed to generate user public id")
	}

	user := &User{
		PublicID: publicID,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("user service: create operation: %w", err)
	}

	s.logger.Info().Msgf("New user registered %s", req.Email)
	return nil
}

func (s *UserService) Login(ctx context.Context, req *LoginUserRequest) (*LoginUserDTO, error) {
	s.logger.Debug().Ctx(ctx).Str("email", req.Email).Msg("Attempt to login")

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errs.NewValidationError("email or password is wrong", err)
	}

	passwordMatch := s.accountSecurity.CheckPasswordHash(ctx, req.Password, user.Password)
	if !passwordMatch {
		return nil, errs.NewValidationError("email or password is wrong", err)
	}

	token, err := s.tokenGenerator.GenerateToken(user.PublicID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	response := &LoginUserDTO{
		PublicID: user.PublicID,
		Token:    token,
	}

	s.logger.Info().Ctx(ctx).Str("userId", response.PublicID).Msg("User login")
	return response, nil
}

func (s *UserService) FindUserById(ctx context.Context, id int) (*UserDTO, error) {
	s.logger.Debug().Ctx(ctx).Int("user id", id).Msg("Attempt to login")
	user, err := s.repo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, errs.NewErrNotFound("user")
		}
		return nil, errs.ErrInternal
	}

	return s.toUserDTO(user), nil
}

func (s *UserService) FindUserByPublicID(ctx context.Context, publicID string) (*UserDTO, error) {
	user, err := s.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, errs.NewErrNotFound("user")
		}
		return nil, errs.ErrInternal
	}
	return s.toUserDTO(user), nil
}

func (s *UserService) GetUserPrimaryID(ctx context.Context, publicID string) (*uint, error) {
	userID, err := s.repo.GetUserPrimaryID(ctx, publicID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, errs.NewErrNotFound("user")
		}
		return nil, errs.ErrInternal
	}
	return userID, nil
}

func (s *UserService) toUserDTO(user *User) *UserDTO {
	return &UserDTO{
		ID:       int(user.ID),
		Email:    user.Email,
		PublicID: user.PublicID,
		Role:     user.Role,
	}
}