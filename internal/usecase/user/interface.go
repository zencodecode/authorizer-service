package user

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type CreateUseCase interface {
	Execute(ctx context.Context, params CreateParams) (*entity.User, error)
}
