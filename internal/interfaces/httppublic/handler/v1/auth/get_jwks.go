package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

func (h *Handler) GetJWKS(c *gin.Context) {
	jwks := auth.BuildJWKS(h.cfg.Auth.JWT.PublicKey, h.cfg.Auth.JWT.KeyID)
	response.Success(c, "get jwks succceed", jwks)
}
