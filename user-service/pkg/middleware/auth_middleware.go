package middleware

import (
	"net/http"
	"strings"

	"github.com/anrisys/quicket/user-service/pkg/errs"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			resp := errs.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "authorization header missing",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			resp := errs.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: "invalid authorization header format",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid {
			resp := errs.ErrorResponse{
				Code: "UNAUTHORIZED",
				Message: "invalid token",
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, resp)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("publicID", claims["sub"])
			c.Set("role", claims["role"])
		}
		c.Next()
	}
}