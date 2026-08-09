package application

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, application *entity.Application) error
	GetByID(ctx context.Context, id string) (*entity.Application, error)
	Update(ctx context.Context, application *entity.Application) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Application, error)
}
