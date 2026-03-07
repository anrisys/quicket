package middleware

import (
	"errors"
	"net/http"
	"quicket/booking-service/internal/booking"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextUserPublicID = "userPublicID"
	ContextUserRole     = "userRole"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(booking.NewAppError(
				http.StatusUnauthorized,
				booking.CodeUnauthorized,
				"authorization header missing",
				nil,
			))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Error(booking.NewAppError(
				http.StatusUnauthorized,
				booking.CodeUnauthorized,
				"invalid authorization header format",
				nil,
			))
			c.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			c.Error(booking.NewAppError(
				http.StatusUnauthorized,
				booking.CodeUnauthorized,
				"invalid token",
				err,
			))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.Error(booking.NewAppError(
				http.StatusUnauthorized,
				booking.CodeUnauthorized,
				"invalid token claims",
				nil,
			))
			c.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			c.Error(booking.NewAppError(
				http.StatusUnauthorized,
				booking.CodeUnauthorized,
				"invalid token subject",
				nil,
			))
			c.Abort()
			return
		}

		role, _ := claims["role"].(string)

		c.Set(ContextUserPublicID, sub)
		c.Set(ContextUserRole, role)

		c.Next()
	}
}
