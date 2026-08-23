package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/interfaces/httppublic/serializer"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
	"github.com/zencodecode/authorizer-service/pkg/validation"
)

type LoginRequest struct {
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
	ChallengeID string `form:"login_challenge"`
}

type LoginResponse struct {
	User              serializer.User `json:"user"`
	AuthorizationCode *string         `json:"authorization_code"`
	ExpiresAt         *int64          `json:"expires_at"`
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

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

	params := toLoginParams(req)

	output, err := h.loginUC.Execute(c.Request.Context(), params)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	r := LoginResponse{
		User:              serializer.SerializeToUser(*output.User),
		AuthorizationCode: output.AuthorizationCode,
		ExpiresAt:         output.ExpiresAt,
	}

	if output.AuthorizationCode != nil {
		response.Success(c, "login success", r)
		return
	}

	c.Redirect(http.StatusFound, "/consent?login_challenge="+output.Session.CodeChallenge)
}

func toLoginParams(r LoginRequest) auth.LoginParams {
	return auth.LoginParams{
		Email:       r.Email,
		Password:    r.Password,
		ChallengeID: r.ChallengeID,
	}
}
