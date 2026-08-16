package organizationapplication

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, orgApp *entity.OrganizationApplication) error
	GetByOrganizationAndApplication(ctx context.Context, organizationID, applicationID uuid.UUID) (*entity.OrganizationApplication, error)
	SetActive(ctx context.Context, organizationID, applicationID uuid.UUID, isActive bool) error
	Delete(ctx context.Context, organizationID, applicationID uuid.UUID) error
	ListApplicationsByOrganization(ctx context.Context, organizationID uuid.UUID) ([]*entity.OrganizationApplication, error)
	ListOrganizationsByApplication(ctx context.Context, applicationID uuid.UUID, limit, offset int) ([]*entity.OrganizationApplication, error)
}
