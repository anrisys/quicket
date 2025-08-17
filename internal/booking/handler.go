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

func (h *Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	userPublicID := c.GetString("publicID")

	var req *dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fieldsErrors := errs.ExtractValidationErrors(err)
		validateErr := errs.NewValidationError("Invalid request data", fieldsErrors, err)
		c.Error(validateErr)
		return
	}

	booking, err := h.svc.Create(ctx, req, userPublicID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": "success",
		"message": "Booking created successfully",
		"booking": booking,
	})
}