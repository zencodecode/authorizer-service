package invitation

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, invitation *entity.Invitation) error
	GetByID(ctx context.Context, id string) (*entity.Invitation, error)
	GetByToken(ctx context.Context, token string) (*entity.Invitation, error)
	UpdateStatus(ctx context.Context, id, status string) error
	Delete(ctx context.Context, id string) error
	ListByOrganization(ctx context.Context, organizationID string, status *string, limit, offset int) ([]*entity.Invitation, error)
}
