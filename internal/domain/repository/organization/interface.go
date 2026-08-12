package organization

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, org *entity.Organization) error
	GetByID(ctx context.Context, id string) (*entity.Organization, error)
	GetBySlug(ctx context.Context, slug string) (*entity.Organization, error)
	Update(ctx context.Context, org *entity.Organization) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]*entity.Organization, error)
}
