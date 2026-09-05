package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
	"github.com/zencodecode/authorizer-service/pkg/validation"
)

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
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

	params := auth.ResetPasswordParams{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}

	result, err := h.resetPasswordUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, result.Message, result)
}
