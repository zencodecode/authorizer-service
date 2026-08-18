package auth

import "context"

type LoginUsecase interface {
	Execute(ctx context.Context, params LoginParams) (*LoginOutput, error)
}

type RegisterUsecase interface {
	Execute(ctx context.Context, params RegisterParams) (*RegisterOutput, error)
}

type AuthorizeUsecase interface {
	Execute(ctx context.Context, params AuthorizeParams) (*AuthorizeResult, error)
}
