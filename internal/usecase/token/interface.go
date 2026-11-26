package token

import "context"

type (
	ExchangeUsecase interface {
		Execute(ctx context.Context, params ExchangeParams) (*ExchangeResult, error)
	}
	RevokeUsecase interface {
		Execute(ctx context.Context, params RevokeParams) (*RevokeResult, error)
	}
	RefreshUsecase interface {
		Execute(ctx context.Context, params RefreshParams) (*RefreshResult, error)
	}
)
