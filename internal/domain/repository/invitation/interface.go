package invitation

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, invitation *entity.Invitation) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Invitation, error)
	GetByToken(ctx context.Context, token string) (*entity.Invitation, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByOrganization(ctx context.Context, organizationID uuid.UUID, status *string, limit, offset int) ([]*entity.Invitation, error)
}
