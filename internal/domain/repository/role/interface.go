package role

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error)
	GetBySlug(ctx context.Context, organizationID *uuid.UUID, applicationID uuid.UUID, slug string) (*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByApplication(ctx context.Context, applicationID uuid.UUID, limit, offset int) ([]*entity.Role, error)
	ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entity.Role, error)
	ListByOrganizationAndApplication(ctx context.Context, organizationID, applicationID uuid.UUID) ([]*entity.Role, error)
}
