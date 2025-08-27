package user

import (
	"context"
	"fmt"

	commonDTO "github.com/anrisys/quicket/internal/dto"
	userDTO "github.com/anrisys/quicket/internal/user/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/anrisys/quicket/pkg/security"
	"github.com/anrisys/quicket/pkg/token"
	"github.com/rs/zerolog"
)

type UserServiceInterface interface {
	Register(ctx context.Context, req *userDTO.RegisterUserRequest) error
	Login(ctx context.Context, req *userDTO.LoginUserRequest) (*userDTO.LoginUserDTO, error)
	FindUserById(ctx context.Context, id int) (*userDTO.UserDTO, error)
}

type UserService struct {
	repo 			UserRepositoryInterface
	logger  		zerolog.Logger
	accountSecurity security.AccountSecurityInterface
	tokenGenerator 	token.GeneratorInterface
}

func NewUserService(
	repo UserRepositoryInterface, 
	logger zerolog.Logger, 
	accountSecurity security.AccountSecurityInterface,
	tokenGenerator 	token.GeneratorInterface,
) *UserService {
	return &UserService{
		repo: repo,
		logger: logger,
		accountSecurity: accountSecurity,
		tokenGenerator: tokenGenerator,
	}
}

func (s *UserService) Register(ctx context.Context, req *userDTO.RegisterUserRequest) error {
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

func (s *UserService) Login(ctx context.Context, req *userDTO.LoginUserRequest) (*userDTO.LoginUserDTO, error) {
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

	response := &userDTO.LoginUserDTO{
		PublicID: user.PublicID,
		Token: token,
	}

	s.logger.Info().Ctx(ctx).Str("userId", response.PublicID).Msg("User login")
	return response, nil
}

func (s *UserService) FindUserById(ctx context.Context, id int) (*userDTO.UserDTO, error) {
	s.logger.Debug().Ctx(ctx).Int("user id", id).Msg("Attempt to login")

	user, err := s.repo.FindById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("user service: findById  %w", err)
	}

	return s.toUserDTO(user), nil
}

func (s *UserService) FindUserByPublicID(ctx context.Context, publicID string) (*commonDTO.UserDTO, error) {
	user, err := s.repo.FindByPublicID(ctx, publicID)
	if err != nil {
		return nil, fmt.Errorf("user service#findByPublicID: %w ", err)
	}
	return &commonDTO.UserDTO{
		ID: int(user.ID),
		PublicID: user.PublicID,
		Email: user.Email,
		Role: user.Role,
	}, nil
}

func (s *UserService) GetUserID(ctx context.Context, publicID string) (*int, error) {
	return s.repo.GetUserID(ctx, publicID)
}

func (s *UserService) toUserDTO(user *User) *userDTO.UserDTO {
	return &userDTO.UserDTO{
		ID: int(user.ID),
		Email: user.Email,
		PublicID: user.PublicID,
		Role: user.Role,
	}
}
