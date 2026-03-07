package middleware

import (
	"errors"
	"net/http"
	"quicket/booking-service/internal/booking"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Recover panic
		defer func() {
			if r := recover(); r != nil {
				log.Error().
					Interface("panic", r).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("panic recovered")

				c.JSON(http.StatusInternalServerError, gin.H{
					"code":    "INTERNAL_ERROR",
					"message": "An unexpected internal server error occurred.",
				})
				c.Abort()
			}
		}()

		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		// Validation Error
		var validationErr *booking.ValidationError
		if errors.As(err, &validationErr) {
			log.Warn().
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Interface("fields", validationErr.Fields).
				Msg("validation error")

			c.JSON(http.StatusBadRequest, gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "Invalid input data",
				"fields":  validationErr.Fields,
			})
			c.Abort()
			return
		}

		// Application Error
		var appErr *booking.AppError
		if errors.As(err, &appErr) {
			log.Error().
				Err(appErr.Err).
				Str("code", appErr.Code).
				Str("path", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Msg("application error")

			c.JSON(appErr.Status, gin.H{
				"code":    appErr.Code,
				"message": appErr.Message,
			})
			c.Abort()
			return
		}

		// Unknown Error
		log.Error().
			Err(err).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Msg("unhandled error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected internal server error occurred.",
		})
		c.Abort()
	}
}
