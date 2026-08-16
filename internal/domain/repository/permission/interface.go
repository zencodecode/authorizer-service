package permission

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, permission *entity.Permission) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error)
	GetBySlug(ctx context.Context, applicationID uuid.UUID, slug string) (*entity.Permission, error)
	Update(ctx context.Context, permission *entity.Permission) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByApplication(ctx context.Context, applicationID uuid.UUID, limit, offset int) ([]*entity.Permission, error)
}
