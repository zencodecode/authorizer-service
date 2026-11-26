package organizationuser

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, orgUser *entity.OrganizationUser) error
	GetByOrganizationAndUser(ctx context.Context, organizationID, userID uuid.UUID) (*entity.OrganizationUser, error)
	UpdateStatus(ctx context.Context, organizationID, userID uuid.UUID, status string) error
	Delete(ctx context.Context, organizationID, userID uuid.UUID) error
	ListUsersByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entity.OrganizationUser, error)
	ListOrganizationsByUser(ctx context.Context, userID uuid.UUID) ([]*entity.OrganizationUser, error)
}
