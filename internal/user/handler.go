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
	c.JSON(http.StatusCreated, gin.H{
		"code": "SUCCESS",
		"message": "User registered successfully",
	})
}

func (h *UserHandler) Login(c *gin.Context)  {
	ctx := c.Request.Context()

	var req dto.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fieldsErrors := errs.ExtractValidationErrors(err)
		validateErr := errs.NewValidationError("Invalid login data", fieldsErrors, err)
		c.Error(validateErr)
		return
	}

	user, err := h.service.Login(ctx, &req)
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

	c.JSON(http.StatusOK, gin.H{
		"code": "SUCCESS",
		"message": "User login successfully",
		"user": user,
	})
}