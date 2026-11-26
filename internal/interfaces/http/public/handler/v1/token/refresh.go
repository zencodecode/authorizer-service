package token

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	tokenuc "github.com/zencodecode/authorizer-service/internal/usecase/token"
	"github.com/zencodecode/authorizer-service/pkg/response"
	"github.com/zencodecode/authorizer-service/pkg/validation"
)

type (
	RefreshRequest struct {
		GrantType    string `form:"grant_type" binding:"required"`
		RefreshToken string `form:"refresh_token" binding:"required"`
		ClientID     string `form:"client_id" binding:"required"`
		ClientSecret string `form:"client_secret" binding:"required"`
	}

	RefreshResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token,omitempty"`
		Scope        string `json:"scope,omitempty"`
		IDToken      string `json:"id_token,omitempty"`
	}
)

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBind(&req); err != nil {
		_ = c.Error(apperr.NewDirectError(enum.INVALID_REQUEST, err.Error()))
		return
	}

	validator := validation.CustomValidator{
		Validator: validation.InitValidator(),
	}

	if err := validator.Validate(c, req); err != nil {
		_ = c.Error(apperr.NewDirectError(enum.INVALID_REQUEST, err.Error()))
		return
	}

	params := toRefreshParams(req)

	result, err := h.refreshUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res := RefreshResponse{
		AccessToken:  result.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		RefreshToken: result.RefreshToken,
		Scope:        result.Scope,
		IDToken:      result.IDToken,
	}
	response.Success(c, "success", res)
}

func toRefreshParams(r RefreshRequest) tokenuc.RefreshParams {
	return tokenuc.RefreshParams{
		GrantType:    r.GrantType,
		RefreshToken: r.RefreshToken,
		ClientID:     r.ClientID,
		ClientSecret: r.ClientSecret,
	}
}
