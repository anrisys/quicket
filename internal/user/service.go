package user

import (
	"context"
	"fmt"

	"github.com/anrisys/quicket/internal/user/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	Register (ctx context.Context, req *dto.RegisterUserRequest) error
	Login (ctx context.Context, req *dto.LoginUserRequest) (*dto.UserDTO, error)
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

func (s *UserService) Register(ctx context.Context, req *dto.RegisterUserRequest) error {
	s.logger.Debug().Str("email", req.Email).Msg("Attempt to register a user")
	
	s.logger.Debug().Msgf("Checking if emai exists %s", req.Email)
	emailRegistered := s.repo.EmailExists(ctx, req.Email)

	if emailRegistered {
		return errs.NewConflictError("email already registered")
	}
	
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return errs.NewInternalError("failed to hash password")
	}

	publicID, err := GeneratePublicID()
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

	passwordMatch := CheckPasswordHash(req.Password, user.Password)
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GeneratePublicID() (string, error) {
	public_id, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	return public_id.String(), nil
}
