package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/interfaces/http/serializer"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
	"github.com/zencodecode/authorizer-service/pkg/validation"
)

type (
	LoginRequest struct {
		Email       string `json:"email" validate:"required"`
		Password    string `json:"password" validate:"required"`
		ChallengeID string `form:"login_challenge"`
	}

	LoginResponse struct {
		User          serializer.User           `json:"user"`
		Organizations []serializer.Organization `json:"organizations"`
		ChallengeID   string                    `json:"login_challenge"`
	}
)

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(apperr.NewDirectError(enum.INVALID_REQUEST, err.Error()))
		return
	}

	if err := c.ShouldBindQuery(&req); err != nil {
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

	params := toLoginParams(req)

	result, err := h.loginUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if result.NextStep == auth.CONSENT {
		organizations := make([]serializer.Organization, 0, len(result.Organizations))
		for _, org := range result.Organizations {
			organizations = append(organizations, serializer.SerializeToOrganization(*org))
		}

		res := LoginResponse{
			User:          serializer.SerializeToUser(*result.User),
			Organizations: organizations,
			ChallengeID:   result.ChallengeID,
		}

		response.Success(c, "proceed to consent", res)
		return
	}
	c.Redirect(http.StatusFound, *result.RedirectURL)
}

func toLoginParams(r LoginRequest) auth.LoginParams {
	return auth.LoginParams{
		Email:       r.Email,
		Password:    r.Password,
		ChallengeID: r.ChallengeID,
	}
}
