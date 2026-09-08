package discovery

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

func (h *Handler) JWKS(c *gin.Context) {
	jwks := auth.BuildJWKS(h.publicKey, h.keyID)
	response.Success(c, "get jwks succceed", jwks)
}
