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
	ExchangeRequest struct {
		GrantType    string `form:"grant_type" binding:"required"`
		Code         string `form:"code" binding:"required"`
		RedirectURI  string `form:"redirect_uri" binding:"required"`
		ClientID     string `form:"client_id" binding:"required"`
		ClientSecret string `form:"client_secret" binding:"required"`
		CodeVerifier string `form:"code_verifier" binding:"required"`
	}

	// ExchangeResponse follows RFC 6749 §5.1 — Successful Response
	ExchangeResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token,omitempty"`
		Scope        string `json:"scope,omitempty"`
		IDToken      string `json:"id_token,omitempty"`
	}
)

func (h *Handler) Exchange(c *gin.Context) {
	var req ExchangeRequest
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

	params := toExchangeParams(req)

	output, err := h.exchangeUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res := ExchangeResponse{
		AccessToken:  output.AccessToken,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		RefreshToken: output.RefreshToken,
		Scope:        "openid profile",
		IDToken:      output.IDToken,
	}
	response.Success(c, "success", res)
}

func toExchangeParams(r ExchangeRequest) tokenuc.ExchangeParams {
	return tokenuc.ExchangeParams{
		GrantType:    r.GrantType,
		Code:         r.Code,
		RedirectURI:  r.RedirectURI,
		ClientID:     r.ClientID,
		ClientSecret: r.ClientSecret,
		CodeVerifier: r.CodeVerifier,
	}
}
