package auth

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
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

type AuthorizeHandler struct {
	authorizeUC auth.AuthorizeUsecase
	logger      service.Logger
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

	_ = output
}

func (h *Handler) handleAuthorizeError(c *gin.Context, params auth.AuthorizeParams, err error) {
	if err != nil {
		var authErr *auth.AuthorizeError

		switch {
		case errors.As(err, &authErr):
			redirectWithError(c, params.RedirectURI, params.State, authErr.Code, authErr.Description)

		case errors.Is(err, auth.ErrInvalidClient), errors.Is(err, auth.ErrRedirectURINotRegistered):
			renderErrorPage(c, 500, err.Error())

		default:
			renderErrorPage(c, 500, "internal error")
		}
		return
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

func redirectWithError(c *gin.Context, redirectURI, state, code, description string) {
	u, _ := url.Parse(redirectURI)
	q := u.Query()
	q.Set("error", code)
	if description != "" {
		q.Set("error_description", description)
	}
	if state != "" {
		q.Set("state", state)
	}
	u.RawQuery = q.Encode()
	http.Redirect(c.Writer, c.Request, u.String(), http.StatusFound)
}

func renderErrorPage(c *gin.Context, status int, message string) {
	c.HTML(status, "error.html", gin.H{
		"Code":    status,
		"Message": message,
	})
}
