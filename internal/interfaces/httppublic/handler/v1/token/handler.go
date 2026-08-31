package token

import (
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/usecase/token"
)

type Handler struct {
	exchangeUC token.ExchangeUsecase
	revokeUC   token.RevokeUsecase
	logger     service.Logger
}

func New(
	exchangeUC token.ExchangeUsecase,
	revokeUC token.RevokeUsecase,
	logger service.Logger,
) *Handler {
	return &Handler{
		exchangeUC: exchangeUC,
		revokeUC:   revokeUC,
		logger:     logger,
	}
}
