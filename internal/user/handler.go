package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	service *ServiceInterface
	logger zerolog.Logger
}

func NewHandler(service *ServiceInterface, logger zerolog.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) Register(c *gin.Context) {
	c.String(http.StatusOK, "Hi this is from Register")
}