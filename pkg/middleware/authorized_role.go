package middleware

import "github.com/gin-gonic/gin"

func AuthorizedRole(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenRole := c.GetString("role")
		for _, role := range allowedRoles {
			if tokenRole == role {
				c.Next()
			}
		}

		c.AbortWithStatusJSON(403, gin.H{
			"error": "insufficient permission",
		})
	}
}