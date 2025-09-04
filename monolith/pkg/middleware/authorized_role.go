package middleware

import (
	"net/http"
	"slices"

	"github.com/anrisys/quicket/pkg/errs"
	"github.com/gin-gonic/gin"
)

func AuthorizedRole(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoleInterface, exists := c.Get("role")
		if !exists {
			resp := errs.ErrorResponse{
				Code: "FORBIDDEN",
				Message: "role information missing",
			}
			c.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}

		userRole, ok := userRoleInterface.(string)
		if !ok {
			resp := errs.ErrorResponse{
				Code: "FORBIDDEN",
				Message: "Invalid role format",
			}
			c.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}

		hasAccess := slices.Contains(allowedRoles, userRole)

		if !hasAccess {
			resp := errs.ErrorResponse{
				Code: "FORBIDDEN",
				Message: "Insufficient permissions",
			}		
			c.AbortWithStatusJSON(http.StatusForbidden, resp)
			return
		}

		c.Next()
	}
}