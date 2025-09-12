package middleware

import (
	"errors"
	"quicket/booking-service/pkg/errs"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error().Stack().Interface("panic_value", r).
					Str("path", c.Request.URL.Path).
					Str("method", c.Request.Method).
					Msg("Recovered from panic in request handler")
				appErr := errs.ErrInternal
				c.JSON(appErr.Status, gin.H{
					"code": appErr.Code,
					"message": appErr.Message,
				})
				c.Abort()
			}
		}()
		
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		var (
			appErr        *errs.AppError
			validationErr *errs.ValidationError
			status        = errs.ErrInternal.Status
			code          = errs.ErrInternal.Code
			message       = errs.ErrInternal.Message
			fields        []errs.FieldError
		)

		logEvent := log.Error().
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method)

		switch {
			case errors.As(err, &validationErr):
				appErr = &validationErr.AppError
				fields = validationErr.Fields
				logEvent.Err(validationErr.Unwrap()).
					Str("type", "validation").
					Interface("fields", fields).
					Msg("Validation failed")

			case errors.As(err, &appErr):
				logEvent.Err(appErr.Unwrap()).
					Str("type", "app").
					Str("code", appErr.Code).
					Int("status", appErr.Status).
					Msg("Application error")

			default:
				logEvent.Err(err).
					Str("type", "unknown").
					Msg("Unhandled error")
		}

		if appErr != nil {
			status = appErr.Status
			code = appErr.Code
			message = appErr.Message
		}

		resp := errs.ErrorResponse{
			Code: code,
			Message: message,
			Fields: fields,
		}

		c.JSON(status, resp)
		c.Abort()
	}
}