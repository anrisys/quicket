package internal

import (
	"net/http"

	"github.com/anrisys/quicket/event-service/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type EventHandler struct {
	EventService EventServiceInterface
	logger       zerolog.Logger
}

func NewEventHandler(eventService EventServiceInterface, logger zerolog.Logger) *EventHandler {
	return &EventHandler{
		EventService: eventService,
		logger:       logger,
	}
}

// Create godoc
// @Summary Create new event
// @Description Create a new event (Admin/Organizer only)
// @Tags Events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateEventRequest true "Event creation data"
// @Success 201 {object} CreateEventSuccessResponse
// @Failure 400 {object} errs.ErrorResponse "Validation error"
// @Failure 401 {object} errs.ErrorResponse "Unauthorized"
// @Failure 403 {object} errs.ErrorResponse "Forbidden (role restriction)"
// @Router /api/v1/events [post]
func (h *EventHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	publicID := c.GetString("publicID")

	var req *CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		validationErr := errs.NewValidationError("Invalid login data", err)
		c.Error(validationErr)
		return
	}

	event, err := h.EventService.Create(ctx, req, publicID)
	if err != nil {
		c.Error(err)
		return
	}

	response := CreateEventSuccessResponse{
		ResponseSuccess: ResponseSuccess{
			Code:    "SUCCESS",
			Message: "Event created successfully",
		},
		Event: SimpleEventDTO{
			PublicID:  event.PublicID,
			Title:     event.Title,
			StartDate: event.StartDate,
			EndDate:   event.EndDate,
		},
	}

	c.JSON(http.StatusCreated, response)
}
