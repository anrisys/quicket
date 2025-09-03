package internal

import (
	"net/http"
	"strconv"

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

// Get userID
// @Summary Retrieve user primary id
// @Description Retrieve user's primary id from user's public id
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param publicID path string true "User Public ID"
// @Success 200 {object} GetPrimaryIDSuccess
// @Failure 401 {object} errs.ErrorResponse "Unauthorized"
// @Failure 404 {object} errs.ErrorResponse "User not found"
// @Failure 500 {object} errs.ErrorResponse "Internal server error"
// @Router /api/v1/users/{publicID}/primary-id [get]
func (h *UserHandler) GetUserPrimaryID(c *gin.Context)  {
	publicID := c.Param("publicID")

	primaryID, err := h.srv.GetUserPrimaryID(c.Request.Context(), publicID)
	if err != nil {
		c.Error(err)
		return
	}
	response := GetPrimaryIDSuccess{
		ResponseSuccess: ResponseSuccess{
			Code: "SUCCESS",
			Message: "Retrieve user's primary id successful",
		},
		PrimaryID: *primaryID,
	}

	c.JSON(http.StatusOK, response)
}

// Get user
// @Summary Retrieve user by ID
// @Description Retrieve user's data from user's id
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param publicID path string true "User ID"
// @Success 200 {object} GetUserByIDSuccess
// @Failure 400 {object} errs.ErrorResponse "Validation error"
// @Failure 401 {object} errs.ErrorResponse "Unauthorized"
// @Failure 404 {object} errs.ErrorResponse "User not found"
// @Failure 500 {object} errs.ErrorResponse "Internal server error"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context)  {
	idParam := c.Param("id")

    id, err := strconv.Atoi(idParam)
    if err != nil {
		response := errs.ErrorResponse{
			Code: "VALIDATION_ERROR",
			Message: "Invalid user ID format. User id must be a valid integer",
		}
        c.JSON(http.StatusBadRequest, response)
        return
    }
    
    if id <= 0 {
		response := errs.ErrorResponse{
			Code: "VALIDATION_ERROR",
			Message: "Invalid user ID format. User id must be positive integer",
		}
        c.JSON(http.StatusBadRequest, response)
        return
    }

	user, err := h.srv.FindUserById(c.Request.Context(), int(id))
	if err != nil {
		c.Error(err)
		return
	}
	response := GetUserByIDSuccess{
		ResponseSuccess: ResponseSuccess{
			Code: "SUCCESS",
			Message: "Get user by id successful",
		},
		Data: *user,
	}

	c.JSON(http.StatusOK, response)
}

// Get user
// @Summary Retrieve user by public ID
// @Description Retrieve user's data from user's public id
// @Tags User
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param publicID path string true "User Public ID"
// @Success 200 {object} GetUserByPublicIDSuccess
// @Failure 401 {object} errs.ErrorResponse "Unauthorized"
// @Failure 404 {object} errs.ErrorResponse "User not found"
// @Failure 500 {object} errs.ErrorResponse "Internal server error"
// @Router /api/v1/users/public/{publicID} [get]
func (h *UserHandler) GetUserByPublicID(c *gin.Context)  {
	publicID := c.Param("publicID")

	user, err := h.srv.FindUserByPublicID(c.Request.Context(), publicID)
	if err != nil {
		c.Error(err)
		return
	}
	response := GetUserByPublicIDSuccess{
		ResponseSuccess: ResponseSuccess{
			Code: "SUCCESS",
			Message: "Get user by public id successful",
		},
		Data: *user,
	}

	c.JSON(http.StatusOK, response)
}