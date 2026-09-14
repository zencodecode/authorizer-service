package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
)

type AuthorizeRequest struct {
	ResponseType        string `form:"response_type" binding:"required"`
	ClientID            string `form:"client_id" binding:"required"`
	RedirectURI         string `form:"redirect_uri" binding:"required"`
	Scope               string `form:"scope" binding:"required"`
	State               string `form:"state" binding:"required"`
	CodeChallenge       string `form:"code_challenge" binding:"required"`
	CodeChallengeMethod string `form:"code_challenge_method" binding:"required"`
}

// Authorize godoc
// @Summary Authorize user
// @Description Generate authorization code
// @Tags Authorize
// @Accept json
// @Produce json
// @Security BearerAccessToken
// @Param request formData AuthorizeRequest true "Authorize payload"
// @Router /oauth2/authorize [get]
func (h *Handler) Authorize(c *gin.Context) {
	var req AuthorizeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		_ = c.Error(apperr.NewDirectError(enum.INVALID_REQUEST, err.Error()))
		return
	}

	params := toAuthorizeParams(req)

	result, err := h.authorizeUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Redirect(http.StatusFound, "/login?login_challenge="+result.LoginChallengeID)
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
