package auth

import "context"

type (
	LoginUsecase interface {
		Execute(ctx context.Context, params LoginParams) (*LoginResult, error)
	}
	RegisterUsecase interface {
		Execute(ctx context.Context, params RegisterParams) (*RegisterResult, error)
	}
	AuthorizeUsecase interface {
		Execute(ctx context.Context, params AuthorizeParams) (*AuthorizeResult, error)
	}
	ConsentUsecase interface {
		Execute(ctx context.Context, params ConsentParams) (*ConsentResult, error)
	}
)
