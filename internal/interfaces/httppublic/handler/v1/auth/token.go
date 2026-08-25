package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/interfaces/httppublic/serializer"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
	"github.com/zencodecode/authorizer-service/pkg/validation"
)

type (
	TokenRequest struct {
		GrantType    string `form:"grant_type" binding:"required"`
		Code         string `form:"code" binding:"required"`
		RedirectURI  string `form:"redirect_uri" binding:"required"`
		ClientID     string `form:"client_id" binding:"required"`
		ClientSecret string `form:"client_secret" binding:"required"`
		CodeVerifier string `form:"code_verifier" binding:"required"`
	}

	TokenResponse struct {
		User         serializer.User `json:"user"`
		AccessToken  string          `json:"access_token"`
		RefreshToken string          `form:"refresh_token"`
	}
)

func (h *Handler) Token(c *gin.Context) {
	var req TokenRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	validator := validation.CustomValidator{
		Validator: validation.InitValidator(),
	}

	if err := validator.Validate(c, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	params := toTokenParams(req)

	output, err := h.tokenUC.Execute(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	r := TokenResponse{
		User:         serializer.SerializeToUser(*output.User),
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
	}

	response.Success(c, "proceed to consent", r)

}

func toTokenParams(r TokenRequest) auth.TokenParams {
	return auth.TokenParams{
		GrantType:    r.GrantType,
		Code:         r.Code,
		RedirectURI:  r.RedirectURI,
		ClientID:     r.ClientID,
		ClientSecret: r.ClientSecret,
		CodeVerifier: r.CodeVerifier,
	}
}
