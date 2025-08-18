package user

import (
	"errors"
	"net/http"

	"github.com/anrisys/quicket/internal/user/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	service UserServiceInterface
	logger zerolog.Logger
}

func NewUserHandler(service UserServiceInterface, logger zerolog.Logger) *UserHandler {
	return &UserHandler{service: service, logger: logger}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with email and password
// @Tags Public, Users
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserRequest true "User Registration data"
// @Success 201 {object} dto.RegisterUserSuccess
// @Failure 400 {object} errs.ErrorResponse
// @Failure 409 {object} errs.ErrorResponse
// @Router /api/v1/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fieldsErrors := errs.ExtractValidationErrors(err)
		validationErr := errs.NewValidationError("Invalid request data", fieldsErrors, err)
		c.Error(validationErr)
		return
	}
	err := h.service.Register(ctx, &req)

	if err != nil {
		c.Error(err)
		return
	}
	response := dto.RegisterUserSuccess{
		ResponseSuccess: dto.ResponseSuccess{
			Code: "SUCCESS",
			Message: "User registered successful",
		},
	}
	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token
// @Tags Public, Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginUserRequest true "User login data"
// @Success 200 {object} dto.LoginUserSuccess
// @Failure 400 {object} errs.ErrorResponse
// @Failure 401 {object} errs.ErrorResponse
// @Router /api/v1//login [post]
func (h *UserHandler) Login(c *gin.Context)  {
	ctx := c.Request.Context()

	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fieldsErrors := errs.ExtractValidationErrors(err)
		validateErr := errs.NewValidationError("Invalid login data", fieldsErrors, err)
		c.Error(validateErr)
		return
	}

	loginData, err := h.service.Login(ctx, &req)
	if err != nil {
		var appErr *errs.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case "INVALID_DATA", "NOT_FOUND":
				c.JSON(http.StatusBadRequest, gin.H{"error": "email or password is wrong"})
			default:
				c.Error(err)
			}
		}
	}
	response := dto.LoginUserSuccess{
		ResponseSuccess: dto.ResponseSuccess{
			Code: "SUCCESS",
			Message: "User logged in successful",
		},
		Data: *loginData,
	}
	c.JSON(http.StatusOK, response)
}