package booking

import (
	"net/http"
	"quicket/booking-service/internal/booking/dto"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usu *UserUsecase
	adu *AdminUsecase
}

func NewHandler(usu *UserUsecase, adu *AdminUsecase) *Handler {
	return &Handler{
		usu: usu,
		adu: adu,
	}
}

// Create godoc
// @Summary Create new booking
// @Description Create a new booking
// @Tags Bookings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateBookingRequest true "Create booking creation data"
// @Success 201 {object} CreateBookingSuccessResponse
// @Failure 400 {object} errs.ErrorResponse "Validation Error"
// @Failure 401 {object} errs.ErrorResponse "Unauthorized"
// @Failure 409 {object} errs.ErrorResponse "Conflict Error"
// @Failure 500 {object} errs.ErrorResponse "Internal Server Error"
// @Router /api/v1/bookings/ [post]
func (h *Handler) CreateBooking(c *gin.Context) {
	ctx := c.Request.Context()
	userPublicID := c.GetString("publicID")

	var req *dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		valErr := NewAppError(http.StatusBadRequest, CodeValidation, "invalid booking request data", err)
		c.Error(valErr)
		return
	}

	booking, err := h.usu.CreateBooking(ctx, req, userPublicID)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.NewSuccessResponse(booking)

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) CancelBooking(c *gin.Context) {
	ctx := c.Request.Context()
	userPublicID := c.GetString("publicID")

	var req *dto.CancelBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		valErr := NewAppError(http.StatusBadRequest, CodeValidation, "invalid cancel booking request data", err)
		c.Error(valErr)
		return
	}

	booking, err := h.usu.CancelBooking(ctx, req, userPublicID)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.NewSuccessResponse(booking)

	c.JSON(http.StatusCreated, response)
}

// HealthCheck godoc
// @Summary Health Check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Router /api/v1/health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	response := gin.H{
		"status":    "healthy",
		"service":   "user-api",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	}

	c.JSON(http.StatusOK, response)
}

/*
// Create godoc
// @Summary Create new booking
// @Description Create a new booking
// @Tags Bookings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateBookingRequest true "Create booking creation data"
// @Success 201 {object} CreateBookingSuccessResponse
// @Failure 400 {object} errs.ErrorResponse "Validation Error"
// @Failure 401 {object} errs.ErrorResponse "Unauthorized"
// @Failure 409 {object} errs.ErrorResponse "Conflict Error"
// @Failure 500 {object} errs.ErrorResponse "Internal Server Error"
// @Router /api/v1/bookings/ [post]
func (h *Handler) AdminBookingList(c *gin.Context) {
	ctx := c.Request.Context()
	userPublicID := c.GetString("publicID")

	var req *dto.AdminBookingListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErr := errs.NewValidationError("Invalid login data", err)
		c.Error(validationErr)
		return
	}

	booking, err := h.srv.Create(ctx, req, userPublicID)
	if err != nil {
		c.Error(err)
		return
	}

	response := CreateBookingSuccessResponse{
		ResponseSuccess: ResponseSuccess{
			Code:    "SUCCESS",
			Message: "Booking created successfully",
		},
		Data: *booking,
	}

	c.JSON(http.StatusCreated, response)
}
*/
