package middleware

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
)

func OAuthErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		if c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err

		var oauthErr *apperr.OAuthError
		if errors.As(err, &oauthErr) {
			if oauthErr.RedirectURI != "" {
				redirectWithOAuthError(c, oauthErr)
				return
			}

			c.AbortWithStatusJSON(mapOAuthErrorToStatus(oauthErr.Code), gin.H{
				"error":             oauthErr.Code,
				"error_description": oauthErr.Description,
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":             enum.SERVER_ERROR,
			"error_description": "internal server error",
		})
	}
}

func redirectWithOAuthError(c *gin.Context, e *apperr.OAuthError) {
	u, err := url.Parse(e.RedirectURI)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":             enum.SERVER_ERROR,
			"error_description": "invalid redirect_uri",
		})
		return
	}

	q := u.Query()
	q.Set("error", string(e.Code))
	if e.Description != "" {
		q.Set("error_description", e.Description)
	}
	if e.State != "" {
		q.Set("state", e.State)
	}
	u.RawQuery = q.Encode()

	c.Redirect(http.StatusFound, u.String())
	c.Abort()
}

func mapOAuthErrorToStatus(code enum.OAuthError) int {
	switch code {
	case enum.INVALID_CLIENT, enum.UNAUTHORIZED_CLIENT:
		return http.StatusUnauthorized
	case enum.ACCESS_DENIED:
		return http.StatusForbidden
	case enum.SERVER_ERROR:
		return http.StatusInternalServerError
	case enum.TEMPORARILY_UNAVAILABLE:
		return http.StatusServiceUnavailable
	default:
		return http.StatusBadRequest
	}
}
