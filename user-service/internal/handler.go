package internal

import (
	"net/http"

	"github.com/anrisys/quicket/user-service/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	srv    UserServiceInterface
	logger zerolog.Logger
}

func NewUserHandler(srv UserServiceInterface, logger zerolog.Logger) *UserHandler {
	return &UserHandler{
		srv: srv,
		logger: logger,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with email and password
// @Tags Public, Users
// @Accept json
// @Produce json
// @Param request body RegisterUserRequest true "User Registration data"
// @Success 201 {object} RegisterUserSuccess
// @Failure 400 {object} errs.ErrorResponse
// @Failure 409 {object} errs.ErrorResponse
// @Router /api/v1/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErr := errs.NewValidationError("Invalid login data", err)
		c.Error(validationErr)
		return
	}
	err := h.srv.Register(ctx, &req)

	if err != nil {
		c.Error(err)
		return
	}
	response := RegisterUserSuccess{
		ResponseSuccess: ResponseSuccess{
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
// @Param request body LoginUserRequest true "User login data"
// @Success 200 {object} LoginUserSuccess
// @Failure 400 {object} errs.ErrorResponse
// @Failure 401 {object} errs.ErrorResponse
// @Router /api/v1//login [post]
func (h *UserHandler) Login(c *gin.Context)  {
	ctx := c.Request.Context()

	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErr := errs.NewValidationError("Invalid login data", err)
		c.Error(validationErr)
		return
	}

	loginData, err := h.srv.Login(ctx, &req)
	if err != nil {
		c.Error(err)
		return
	}
	response := LoginUserSuccess{
		ResponseSuccess: ResponseSuccess{
			Code: "SUCCESS",
			Message: "User logged in successful",
		},
		Data: *loginData,
	}
	c.JSON(http.StatusOK, response)
}