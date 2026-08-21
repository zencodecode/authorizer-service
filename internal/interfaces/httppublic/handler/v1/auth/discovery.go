package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/auth"
)

func (h *Handler) Discovery(c *gin.Context) {
	openID := auth.BuildOpenIDConfiguration(h.cfg.Auth.OIDC.Issuer)
	c.JSON(http.StatusOK, openID)
}
