package internal

import (
	"net/http"
	"quicket/booking-service/pkg/errs"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	srv ServiceInterface
}

func NewHandler(srv ServiceInterface) *Handler {
	return &Handler{
		srv: srv,
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

	var req *CreateBookingRequest
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
			Code: "SUCCESS",
			Message: "Booking created successfully",
		},
		Data: *booking,
	}

	c.JSON(http.StatusCreated, response)
}