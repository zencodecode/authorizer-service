package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/interfaces/http/public/serializer"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
	"github.com/zencodecode/authorizer-service/pkg/validation"
)

type (
	RegisterRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
		Name     string `json:"name" validate:"required"`
	}
	RegisterResponse struct {
		User serializer.User `json:"user"`
	}
)

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
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

	params := toRegisterParams(req)

	result, err := h.registerUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	res := RegisterResponse{
		User: serializer.User{
			ID:    result.User.ID,
			Email: result.User.Email,
			Name:  result.User.Name,
		},
	}

	response.Success(c, "user registered successfully", res)
}

func toRegisterParams(req RegisterRequest) auth.RegisterParams {
	return auth.RegisterParams{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}
}
