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

func TestUserService_Register(t *testing.T) {
	ctx := context.Background()
	logger := zerolog.Nop()
	validRequest := &dto.RegisterUserRequest{
		Email: "valid@example.com",
		Password: "SecurePass123!",
	}
	
	t.Run("Success - New User Registration", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewUserService(mockRepo, logger)

		mockRepo.On("EmailExists", ctx, "valid@example.com").Return(false)
		mockRepo.On("Create", ctx, mock.MatchedBy(func(u *User) bool  {
			return u.Email == "valid@example.com" &&
				u.Role == "user" && 
				len(u.PublicID) > 0
		})).Return(nil)

		err := service.Register(ctx, validRequest)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Failure - Email Already Exists", func(t *testing.T) {
		mockRepo := new(MockUserRepo)
		service := NewUserService(mockRepo, logger)

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
		assert.Equal(t, "CONFLICT", appErr.Code) 
		mockRepo.AssertNotCalled(t, "Create")
	})
}