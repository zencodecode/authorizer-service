package organizationuser

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, orgUser *entity.OrganizationUser) error
	GetByOrganizationAndUser(ctx context.Context, organizationID, userID string) (*entity.OrganizationUser, error)
	UpdateStatus(ctx context.Context, organizationID, userID, status string) error
	Delete(ctx context.Context, organizationID, userID string) error
	ListUsersByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*entity.OrganizationUser, error)
	ListOrganizationsByUser(ctx context.Context, userID string) ([]*entity.OrganizationUser, error)
}
