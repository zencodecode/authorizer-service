package token

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/usecase/token"
)

type Handler struct {
	exchangeUC token.ExchangeUsecase
	refreshUC  token.RefreshUsecase
	revokeUC   token.RevokeUsecase
	logger     service.Logger
}

func New(
	exchangeUC token.ExchangeUsecase,
	refreshUC token.RefreshUsecase,
	revokeUC token.RevokeUsecase,
	logger service.Logger,
) *Handler {
	return &Handler{
		exchangeUC: exchangeUC,
		refreshUC:  refreshUC,
		revokeUC:   revokeUC,
		logger:     logger,
	}
}

func (h *Handler) Token(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	switch grantType {
	case "authorization_code":
		h.Exchange(c)
	case "refresh_token":
		h.Refresh(c)
	default:
		c.JSON(400, gin.H{
			"error":             "unsupported_grant_type",
			"error_description": "grant_type must be authorization_code or refresh_token",
		})
	}
}
