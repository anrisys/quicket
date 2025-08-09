package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "authorization header missing"})
			return
		}

		// token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		// 	if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok{
		// 		return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		// 	}
		// 	return []byte(secret), nil
		// })
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid{
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// c.Set("jwtClaims", token.Claims)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("publicID", claims["sub"])
			c.Set("role", claims["role"])
		}
		c.Next()
	}
}