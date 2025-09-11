package internal

import (
	"net/http"
	"strconv"

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

// GetEventByPublicID godoc  
// @Summary Get event detail by public ID
// @Description Get event details using public identifier (for external clients)
// @Tags Events - Public
// @Produce json
// @Param public_id path string true "Public Event ID"
// @Success 200 {object} FindByPublicIDSuccessResponse
// @Failure 404 {object} errs.ErrorResponse "Event not found"
// @Failure 500 {object} errs.ErrorResponse "Internal server error"
// @Router /api/v1/events/{public_id} [get]
func (h *EventHandler) GetEventByPublicID(c *gin.Context) {
	ctx := c.Request.Context()
	publicID := c.Param("publicID")

	ev, err := h.EventService.FindByPublicID(ctx, publicID)
	if err != nil {
		c.Error(err)
		return 
	}
	response := FindByPublicIDSuccessResponse{
		ResponseSuccess: ResponseSuccess{
			Code: "SUCCESS",
			Message: "Event successfully retrieved",
		},
		Data: *ev,
	}

	c.JSON(http.StatusOK, response)
}

// GetEventByID godoc
// @Summary Get event detail by internal ID
// @Description Get event details using internal system ID (for internal use only)
// @Tags Events - Internal
// @Produce json
// @Param id path int true "Internal Event ID" 
// @Success 200 {object} FindByIDSuccessResponse
// @Failure 400 {object} errs.ErrorResponse "Invalid ID format"
// @Failure 404 {object} errs.ErrorResponse "Event not found"
// @Failure 500 {object} errs.ErrorResponse "Internal server error"
// @Router /internal/api/v1/events/{id} [get]
func (h *EventHandler) GetEventByID(c *gin.Context) {
	ctx := c.Request.Context()
	idParam := c.Param("id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		errResponse := errs.ErrorResponse{
			Code: "INVALID_PATH_PARAM",
			Message: "Invalid event ID format",
		}
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	if id == 0 {
		errResponse := errs.ErrorResponse{
			Code: "INVALID_PATH_PARAM",
			Message: "ID can not be zero",
		}
		c.JSON(http.StatusBadRequest, errResponse)
	}

	ev, err := h.EventService.FindByID(ctx, uint(id))
	if err != nil {
		c.Error(err)
		return 
	}
	response := FindByIDSuccessResponse{
		ResponseSuccess: ResponseSuccess{
			Code: "SUCCESS",
			Message: "Event successfully retrieved",
		},
		Data: *ev,
	}

	c.JSON(http.StatusOK, response)
}