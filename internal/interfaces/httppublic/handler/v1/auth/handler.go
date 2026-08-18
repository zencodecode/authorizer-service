package auth

import (
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
)

type Handler struct {
	authorizeUC auth.AuthorizeUsecase
	cfg         *config.Config
	logger      service.Logger
}

func New(
	authorizeUC auth.AuthorizeUsecase,
	cfg *config.Config,
	logger service.Logger,
) *Handler {
	return &Handler{
		authorizeUC: authorizeUC,
		cfg:         cfg,
		logger:      logger,
	}
}
