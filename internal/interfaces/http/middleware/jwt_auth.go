package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

func JWTAuth(jwtSvc service.JWTService, log service.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := GetRequestID(c)
		token, err := extractBearerToken(c)

		if err != nil {
			log.Warn(c, "Authentication failed: missing token",
				"request_id", requestID,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"ip", c.ClientIP(),
				"error", err.Error(),
			)

			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		claims, err := jwtSvc.ValidateAccessToken(c, token)
		if err != nil {
			message := mapValidationError(err)

			log.Warn(c, "Authentication failed: token validation error",
				"request_id", requestID,
				"path", c.Request.URL.Path,
				"method", c.Request.Method,
				"ip", c.ClientIP(),
				"error", message,
				"token", truncateToken(token),
			)

			response.Unauthorized(c, message)
			c.Abort()
			return
		}

		SetClaims(c, claims)
		SetToken(c, token)

		log.Info(c, "Authentication successful",
			"request_id", requestID,
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"user_id", claims.Subject,
			"email", claims.Email,
			"token", truncateToken(token),
		)

		c.Next()
	}
}

func extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("missing authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", fmt.Errorf("empty bearer token")
	}

	return token, nil
}

func mapValidationError(err error) string {
	errMsg := err.Error()
	fmt.Println(errMsg)
	switch {
	case strings.Contains(errMsg, "Token expired"):
		return "Token expired"
	case strings.Contains(errMsg, "Invalid token issuer"):
		return "Invalid token issuer"
	case strings.Contains(errMsg, "Invalid token audience"):
		return "Invalid token audience"
	case strings.Contains(errMsg, "failed to get public key"):
		return "Invalid token signature"
	case strings.Contains(errMsg, "unexpected signing method"):
		return "Invalid token signature"
	case strings.Contains(errMsg, "missing or invalid kid"):
		return "Invalid token signature"
	case strings.Contains(errMsg, "failed to fetch JWKS"):
		return "Authentication service unavailable"
	case strings.Contains(errMsg, "JWKS endpoint"):
		return "Authentication service unavailable"
	default:
		return "Invalid token signature"
	}
}

func truncateToken(token string) string {
	if len(token) < 4 {
		return ""
	}
	return "..." + token[len(token)-4:]
}
