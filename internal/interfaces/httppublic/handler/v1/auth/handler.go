package auth

import (
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/auth"
)

type Handler struct {
	jwksSvc auth.JWKSService
	cfg     *config.Config
	logger  service.Logger
}

func New(jwksSvc auth.JWKSService,
	cfg *config.Config,
	logger service.Logger,
) *Handler {
	return &Handler{
		jwksSvc: jwksSvc,
		cfg:     cfg,
		logger:  logger,
	}
}
