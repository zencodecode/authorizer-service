package auth

import (
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
)

type Handler struct {
	authorizeUC auth.AuthorizeUsecase
	loginUC     auth.LoginUsecase
	consentUC   auth.ConsentUsecase
	cfg         *config.Config
	logger      service.Logger
}

func New(
	authorizeUC auth.AuthorizeUsecase,
	loginUC auth.LoginUsecase,
	consentUC auth.ConsentUsecase,
	cfg *config.Config,
	logger service.Logger,
) *Handler {
	return &Handler{
		authorizeUC: authorizeUC,
		loginUC:     loginUC,
		consentUC:   consentUC,
		cfg:         cfg,
		logger:      logger,
	}
}
