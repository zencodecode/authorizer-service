package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/interfaces/httppublic/handler/serializer"
	"github.com/zencodecode/authorizer-service/pkg/response"
)

func (h *Handler) GetJWKS(c *gin.Context) {
	jwks, err := h.jwksSvc.GetJWKS(h.cfg.Auth.JWT.PublicKey, h.cfg.Auth.JWT.KeyID)
	if err != nil {
		h.logger.Error(c, "Failed to generate JWKS",
			"error", err.Error(),
		)
		response.InternalServerError(c, "failed to generate JWKS")
	}

	data := serializer.JWKSResponse{
		Keys: jwks.Keys,
	}

	response.Success(c, "User registered successfully", data)

}
