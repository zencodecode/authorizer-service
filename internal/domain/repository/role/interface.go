package role

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetByID(ctx context.Context, id string) (*entity.Role, error)
	GetBySlug(ctx context.Context, organizationID *string, applicationID, slug string) (*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, id string) error
	ListByApplication(ctx context.Context, applicationID string, limit, offset int) ([]*entity.Role, error)
	ListByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*entity.Role, error)
	ListByOrganizationAndApplication(ctx context.Context, organizationID, applicationID string) ([]*entity.Role, error)
}
