package application

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, app *entity.Application) error
	GetByID(ctx context.Context, id string) (*entity.Application, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Application, error)
	GetByClientID(ctx context.Context, clientID string) (*entity.Application, error)
	Update(ctx context.Context, app *entity.Application) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Application, error)
}
