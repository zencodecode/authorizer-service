package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

func (h *Handler) UserInfo(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		_ = c.Error(apperr.NewDirectError(enum.ACCESS_DENIED, "missing or invalid authorization header"))
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims, err := h.jwtSvc.ValidateAccessToken(c.Request.Context(), tokenString)
	if err != nil {
		h.logger.Error(c.Request.Context(), "invalid access token",
			"action", "LOGOUT",
			"error", err.Error(),
		)
		_ = c.Error(apperr.NewDirectError(enum.ACCESS_DENIED, "invalid or expired access token"))
		return
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		_ = c.Error(apperr.NewDirectError(enum.SERVER_ERROR, "invalid user id in token"))
		return
	}

	if len(claims.Scopes) == 0 {
		_ = c.Error(apperr.NewDirectError(enum.ACCESS_DENIED, "missing scopes in token"))
		return
	}

	params := auth.UserInfoParams{
		UserID: userID,
		Scopes: claims.Scopes,
	}

	res, err := h.userInfoUC.Execute(c.Request.Context(), params)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.Success(c, "user info retrieved successfully", res)
}
