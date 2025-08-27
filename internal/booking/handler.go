package booking

import (
	"net/http"

	"github.com/anrisys/quicket/internal/booking/dto"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	svc ServiceInterface
	logger zerolog.Logger
}

func NewHandler(svc ServiceInterface, logger zerolog.Logger) *Handler {
	return &Handler{
		svc: svc,
		logger: logger,
	}
}

// Create godoc
// @Summary Create new booking
// @Description Create a new booking
// @Tags Bookings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateBookingRequest true "Event creation data" 
// @Success 201 {object} dto.CreateBookingSuccessResponse
// @Failure 400 {object} errs.ErrorResponse "Validation error"
// @Failure 401 {object} errs.ErrorResponse "Unauthorized"
// @Router /api/v1/events/bookings/:eventID [post]
func (h *Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	userPublicID := c.GetString("publicID")
	eventID := c.Param("eventID")

	var req *dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// fieldsErrors := errs.ExtractValidationErrors(err)
		// validateErr := errs.NewValidationError("Invalid login data", fieldsErrors, err)
		validationErr := errs.NewValidationError("Invalid login data", err)
		c.Error(validationErr)
		return
	}

	booking, err := h.svc.Create(ctx, req, userPublicID, eventID)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.CreateBookingSuccessResponse{
		ResponseSuccess: dto.ResponseSuccess{
			Code: "SUCCESS",
			Message: "Booking created successfully",
		},
		Booking: *booking,
	}

	c.JSON(http.StatusCreated, response)
}