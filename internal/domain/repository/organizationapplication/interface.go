package organizationapplication

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, orgApp *entity.OrganizationApplication) error
	GetByOrganizationAndApplication(ctx context.Context, organizationID, applicationID string) (*entity.OrganizationApplication, error)
	SetActive(ctx context.Context, organizationID, applicationID string, isActive bool) error
	Delete(ctx context.Context, organizationID, applicationID string) error
	ListApplicationsByOrganization(ctx context.Context, organizationID string) ([]*entity.OrganizationApplication, error)
	ListOrganizationsByApplication(ctx context.Context, applicationID string, limit, offset int) ([]*entity.OrganizationApplication, error)
}
