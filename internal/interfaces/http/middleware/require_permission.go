package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

func RequirePermission(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := GetClaims(c)
		if claims == nil {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		for _, p := range claims.Permissions {
			if p == requiredPermission {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "insufficient permissions")
		c.Abort()
	}
}

func RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := GetClaims(c)
		if claims == nil {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		permSet := make(map[string]struct{}, len(claims.Permissions))
		for _, p := range claims.Permissions {
			permSet[p] = struct{}{}
		}

		for _, required := range permissions {
			if _, ok := permSet[required]; ok {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "insufficient permissions")
		c.Abort()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := GetClaims(c)
		if claims == nil {
			response.Unauthorized(c, "authentication required")
			c.Abort()
			return
		}

		for _, r := range claims.Roles {
			if r == requiredRole {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "insufficient role")
		c.Abort()
	}
}
