package middleware

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
)

// OAuthErrorHandler is a centralized error handler middleware that inspects
// c.Errors after the handler chain completes and responds according to
// RFC 6749 error semantics:
//   - Redirectable errors (RedirectURI present): 302 redirect with error query params
//   - Fatal errors (no RedirectURI): JSON body with mapped HTTP status code
func OAuthErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		// Prevent double-write if handler already wrote a response
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

			// Fatal: render JSON error, do NOT redirect
			c.AbortWithStatusJSON(mapOAuthErrorToStatus(oauthErr.Code), gin.H{
				"error":             oauthErr.Code,
				"error_description": oauthErr.Description,
			})
			return
		}

		// Fallback for errors that are not *apperr.OAuthError
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":             enum.SERVER_ERROR,
			"error_description": "internal server error",
		})
	}
}

func redirectWithOAuthError(c *gin.Context, e *apperr.OAuthError) {
	u, err := url.Parse(e.RedirectURI)
	if err != nil {
		// redirect_uri is corrupt despite having been validated earlier —
		// fall back to a fatal JSON error rather than redirecting to a broken URL.
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

// mapOAuthErrorToStatus maps OAuth 2.0 error codes to HTTP status codes
// following RFC 6749 §5.2 conventions and common industry practice.
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
		// invalid_request, invalid_grant, unsupported_response_type, invalid_scope
		return http.StatusBadRequest
	}
}
