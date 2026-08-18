package auth

import "context"

type LoginUsecase interface {
	Execute(ctx context.Context, params LoginParams) (*LoginResult, error)
}

type RegisterUsecase interface {
	Execute(ctx context.Context, params RegisterParams) (*RegisterResult, error)
}

type AuthorizeUsecase interface {
	Execute(ctx context.Context, params AuthorizeParams) (*AuthorizeResult, error)
}
