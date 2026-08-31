package auth

import "context"

type (
	AuthorizeUsecase interface {
		Execute(ctx context.Context, params AuthorizeParams) (*AuthorizeResult, error)
	}
	LoginUsecase interface {
		Execute(ctx context.Context, params LoginParams) (*LoginResult, error)
	}
	ConsentUsecase interface {
		Execute(ctx context.Context, params ConsentParams) (*ConsentResult, error)
	}
	LogoutUsecase interface {
		Execute(ctx context.Context, params LogoutParams) error
	}
)
