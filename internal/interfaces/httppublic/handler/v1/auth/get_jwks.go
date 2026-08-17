package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/auth"
)

func (h *Handler) GetJWKS(c *gin.Context) {
	jwks := auth.BuildJWKS(h.cfg.Auth.JWT.PublicKey, h.cfg.Auth.JWT.KeyID)
	c.JSON(http.StatusOK, jwks)
}
