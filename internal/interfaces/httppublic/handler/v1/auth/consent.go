package auth

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
	"github.com/zencodecode/authorizer-service/pkg/validation"
)

type (
	ConsentRequest struct {
		OrganizationID string `json:"organization_id" validate:"required"`
		ChallengeID    string `form:"login_challenge"`
	}
)

func (h *Handler) Consent(c *gin.Context) {
	var req ConsentRequest
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

	params := toConsentParams(req)

	output, err := h.consentUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Redirect(http.StatusFound, *output.RedirectURL)
}

func toConsentParams(r ConsentRequest) auth.ConsentParams {
	OrgID, _ := uuid.Parse(r.OrganizationID)
	return auth.ConsentParams{
		OrgID:       OrgID,
		ChallengeID: r.ChallengeID,
	}
}
