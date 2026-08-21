package auth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

type AuthorizeRequest struct {
	ResponseType        string `form:"response_type" binding:"required"`
	ClientID            string `form:"client_id" binding:"required"`
	RedirectURI         string `form:"redirect_uri" binding:"required"`
	Scope               string `form:"scope"`
	State               string `form:"state"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
}

func (h *Handler) Authorize(c *gin.Context) {
	var req AuthorizeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	params := toAuthorizeParams(req)

	output, err := h.authorizeUC.Execute(c.Request.Context(), params)
	if err != nil {
		h.handleAuthorizeError(c, params, err)
		return
	}

	c.Redirect(http.StatusFound, "/login?login_challenge="+output.LoginChallengeID)
}

func (h *Handler) handleAuthorizeError(c *gin.Context, params auth.AuthorizeParams, err error) {
	var authErr *auth.AuthError

	switch {
	// Redirect-safe errors: redirect back to client with error params
	case errors.As(err, &authErr):
		redirectWithError(c, params.RedirectURI, params.State, authErr.Code, authErr.Description)

	// Non-redirect errors: show error page directly (NEVER redirect to unvalidated URI)
	case errors.Is(err, auth.ErrInvalidClient):
		renderErrorPage(c, http.StatusBadRequest, "Invalid client application")

	case errors.Is(err, auth.ErrRedirectURINotRegistered):
		renderErrorPage(c, http.StatusBadRequest, "Redirect URI is not registered for this application")

	default:
		renderErrorPage(c, http.StatusInternalServerError, "An unexpected error occurred")
	}
}

func toAuthorizeParams(r AuthorizeRequest) auth.AuthorizeParams {
	var scopes []string
	if r.Scope != "" {
		scopes = strings.Fields(r.Scope)
	}

	return auth.AuthorizeParams{
		ResponseType:        r.ResponseType,
		ClientID:            r.ClientID,
		RedirectURI:         r.RedirectURI,
		Scope:               scopes,
		State:               r.State,
		CodeChallenge:       r.CodeChallenge,
		CodeChallengeMethod: r.CodeChallengeMethod,
	}
}

func redirectWithError(c *gin.Context, redirectURI, state, errCode, description string) {
	u, parseErr := url.Parse(redirectURI)
	if parseErr != nil {
		renderErrorPage(c, http.StatusInternalServerError, "Invalid redirect URI")
		return
	}

	q := u.Query()
	q.Set("error", errCode)
	if description != "" {
		q.Set("error_description", description)
	}
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()

	c.Redirect(http.StatusFound, u.String())
}

func renderErrorPage(c *gin.Context, status int, message string) {
	// TODO: Replace with proper HTML template when login UI is implemented
	c.JSON(status, gin.H{
		"error":   http.StatusText(status),
		"message": message,
	})
}
