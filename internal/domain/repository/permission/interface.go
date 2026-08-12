package permission

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, permission *entity.Permission) error
	GetByID(ctx context.Context, id string) (*entity.Permission, error)
	GetBySlug(ctx context.Context, applicationID, slug string) (*entity.Permission, error)
	Update(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, id string) error
	ListByApplication(ctx context.Context, applicationID string, limit, offset int) ([]*entity.Permission, error)
}
