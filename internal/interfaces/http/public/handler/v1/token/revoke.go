package token

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	tokenuc "github.com/zencodecode/authorizer-service/internal/usecase/token"
	"github.com/zencodecode/authorizer-service/pkg/response"
	"github.com/zencodecode/authorizer-service/pkg/validation"
)

type RevokeRequest struct {
	Token         string `form:"token" binding:"required"`
	TokenTypeHint string `form:"token_type_hint" binding:"required"`
	ClientID      string `form:"client_id" binding:"required"`
	ClientSecret  string `form:"client_secret" binding:"required"`
}

func (h *Handler) Revoke(c *gin.Context) {
	var req RevokeRequest
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

	params := toRevokeParams(req)

	result, err := h.revokeUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, result.Message, nil)
}

func toRevokeParams(r RevokeRequest) tokenuc.RevokeParams {
	return tokenuc.RevokeParams{
		Token:         r.Token,
		TokenTypeHint: r.TokenTypeHint,
		ClientID:      r.ClientID,
		ClientSecret:  r.ClientSecret,
	}
}
