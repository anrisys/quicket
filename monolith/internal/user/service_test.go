package user

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/anrisys/quicket/internal/user/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

type MockAccountSecurity struct {
	mock.Mock
}

type MockGenerator struct {
	mock.Mock
}

func (m *MockAccountSecurity) HashPassword(ctx context.Context, password string) (string, error)  {
	args := m.Called(ctx, password)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockAccountSecurity) CheckPasswordHash(ctx context.Context, password, hashedPassword string) bool {
	args := m.Called(ctx, password, hashedPassword)
	return args.Bool(0)
}

func (m *MockAccountSecurity) GeneratePublicID(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) FindById(ctx context.Context, id int) (*User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) EmailExists(ctx context.Context, email string) bool {
	args := m.Called(ctx, email)
	return args.Bool(0) // Special handling for bool return
}

func (m *MockUserRepo) FindByPublicID(ctx context.Context, publicID string) (*User, error) {
	args := m.Called(ctx, publicID)
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepo) GetUserID(ctx context.Context, publicID string) (*uint, error) {
	args := m.Called(ctx, publicID)
	return args.Get(0).(*uint), args.Error(1)
}

func (m *MockGenerator) GenerateToken(publicID, role string) (string, error) {
	args := m.Called(publicID, role)
	return args.Get(0).(string), args.Error(1)
}

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	validRequest := &dto.RegisterUserRequest{
		Email: "valid@example.com",
		Password: "SecurePass123!",
	}
	
	t.Run("Success - New User Registration", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		mockSecurity := new(MockAccountSecurity)
		mockGenerator := new(MockGenerator)
		service := NewUserService(mockRepo, logger, mockSecurity, mockGenerator)

		mockRepo.On("EmailExists", ctx, "valid@example.com").Return(false)
		mockSecurity.On("HashPassword", ctx, validRequest.Password).
			Return("Hashed_" + validRequest.Password, nil).Once()
		mockSecurity.On("GeneratePublicID", ctx).Return("X_PUBLIC_ID", nil).Once()
		mockRepo.On("Create", ctx, mock.MatchedBy(func(u *User) bool  {
			return u.Email == "valid@example.com" &&
				u.Role == "user" && 
				u.Password == "Hashed_" + validRequest.Password &&
				u.PublicID == "X_PUBLIC_ID"
		})).Return(nil)

		err := service.Register(ctx, validRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockSecurity.AssertExpectations(t)
	})

	t.Run("Failure - Email Already Exists", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		mockSecurity := new(MockAccountSecurity)
		mockGenerator := new(MockGenerator)
		service := NewUserService(mockRepo, logger, mockSecurity, mockGenerator)

		mockRepo.On("EmailExists", ctx, "exists@example.com").Return(true)

		err := service.Register(ctx, &dto.RegisterUserRequest{
			Email: "exists@example.com",
			Password: "AnyPassword123!",
		})

		assert.Error(t, err)
		var appErr *errs.AppError
		assert.True(t, errors.As(err, &appErr), "Error should be of type *errs.AppError")
		assert.Equal(t, "email already registered", appErr.Message)
		assert.Equal(t, http.StatusConflict, appErr.Status)
		assert.Equal(t, "CONFLICT_ERROR", appErr.Code) 
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestUserService_Login(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	validRequest := &dto.LoginUserRequest{
		Email: "valid@example.com",
		Password: "SecurePass123!",
	}
	
	t.Run("Success", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		mockSecurity := new(MockAccountSecurity)
		mockGenerator := new(MockGenerator)
		service := NewUserService(mockRepo, logger, mockSecurity, mockGenerator)
		hashedPass := "!@#$hashedPass"
		testUser := &User{
			Email: validRequest.Email,
			Password: hashedPass,
			PublicID: "user_123",
			Role: "user",
		}

		mockRepo.On("FindByEmail", ctx, validRequest.Email).Return(testUser, nil)
		mockSecurity.On("CheckPasswordHash", ctx, validRequest.Password, hashedPass).
			Return(true)

		mockGenerator.On("GenerateToken", testUser.PublicID, testUser.Role).
		Return("mock_token_string", nil)

		userDTO, err := service.Login(ctx, validRequest)

		assert.NoError(t, err)
		assert.Equal(t, &dto.UserDTO{
			Email: validRequest.Email,
			PublicID: "user_123",
			Role: "user",
		}, userDTO)

		mockRepo.AssertExpectations(t)
		mockSecurity.AssertExpectations(t)
	})
}