package middleware

import (
	"time"

	"github.com/anrisys/quicket/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func ZerologLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			id, err := uuid.NewRandom()
			if err != nil {
				c.Error(errs.NewInternalError("failed to generate requestID"))
			}
			requestID = id.String()
			c.Header("X-Request-ID", requestID)
		}

		requestLogger := log.With().
			Str("request_id", requestID).
			Str("path", c.Request.URL.Path).
			Str("method", c.Request.Method).
			Logger()
		
		ctx := requestLogger.WithContext(c.Request.Context())

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		duration := time.Since(start)
		log.Ctx(ctx).Info().
			Int("status", c.Writer.Status()).
			Dur("duration_ms", duration).
			Msg("Request completed")
	}
}