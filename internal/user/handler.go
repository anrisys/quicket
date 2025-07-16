package user

import (
	"net/http"

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
	c.String(http.StatusOK, "Hi this is from Register")
}