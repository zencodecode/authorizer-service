package discovery

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/auth"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

func (h *Handler) OpenIDConfiguration(c *gin.Context) {
	openID := auth.BuildOpenIDConfiguration(h.issuer)
	response.Success(c, "discovery succceed", openID)
}
