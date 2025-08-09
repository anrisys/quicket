package event

import (
	"net/http"

	"github.com/anrisys/quicket/internal/event/dto"
	"github.com/anrisys/quicket/internal/user"
	"github.com/anrisys/quicket/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type EventHandler struct {
	EventService EventServiceInterface
	UserService user.UserDTOServiceInterface
	logger  zerolog.Logger
}

func NewEventHandler(eventService EventServiceInterface, userService user.UserDTOServiceInterface, logger zerolog.Logger) *EventHandler {
	return &EventHandler{
		EventService: eventService,
		UserService: userService,
		logger: logger,
	}
}

func (h *EventHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()
	publicID := c.GetString("publicID")

	userDTO, err := h.UserService.FindUserByPublicID(ctx, publicID)
	if err != nil {
		c.Error(err)
		return
	}

	var req *dto.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fieldsErrors := errs.ExtractValidationErrors(err)
		validateErr := errs.NewValidationError("Invalid request data", fieldsErrors, err)
		c.Error(validateErr)
		return
	}

	event, err := h.EventService.Create(ctx, req, userDTO.ID)
	if err != nil {
		c.Error(err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"code": "success",
		"message": "Event created successfully",
		"event": event,
	})
}
