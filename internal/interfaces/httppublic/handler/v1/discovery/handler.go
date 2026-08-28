package discovery

import (
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type Handler struct {
	cfg    *config.Config
	logger service.Logger
}

func New(
	cfg *config.Config,
	logger service.Logger,
) *Handler {
	return &Handler{
		cfg:    cfg,
		logger: logger,
	}
}
