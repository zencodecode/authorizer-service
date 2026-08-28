package user

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type CreateUseCase interface {
	Execute(ctx context.Context, params CreateParams) (*entity.User, error)
}

type RegisterUsecase interface {
	Execute(ctx context.Context, params RegisterParams) (*RegisterResult, error)
}

type RegisterParams struct {
	Email    string
	Password string
	Name     string
}

type RegisterResult struct {
	User *entity.User
}
