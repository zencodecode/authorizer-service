package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

const (
	RequestIDKey     = "request_id"
	ClaimsContextKey = "auth_claims"
	TokenContextKey  = "auth_token"
)

func SetRequestID(c *gin.Context, requestID string) {
	c.Set(RequestIDKey, requestID)
}

func SetClaims(c *gin.Context, claims *entity.Claims) {
	c.Set(ClaimsContextKey, claims)
}

func SetToken(c *gin.Context, token string) {
	c.Set(TokenContextKey, token)
}

func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

func GetToken(c *gin.Context) string {
	value, exists := c.Get(TokenContextKey)
	if !exists {
		return ""
	}
	token, ok := value.(string)
	if !ok {
		return ""
	}
	return token
}

func GetClaims(c *gin.Context) (*entity.Claims, error) {
	value, exists := c.Get(ClaimsContextKey)
	if !exists {
		return nil, errors.New("claims not found in context")
	}

	claims, ok := value.(*entity.Claims)
	if !ok {
		return nil, errors.New("invalid claims type in context")
	}

	return claims, nil
}

func MustGetClaims(c *gin.Context) *entity.Claims {
	claims, err := GetClaims(c)
	if err != nil {
		panic("claims not found in context: " + err.Error())
	}
	return claims
}
