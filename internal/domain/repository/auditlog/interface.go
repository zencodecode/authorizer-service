package auditlog

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, log *entity.AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditLog, error)
	ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entity.AuditLog, error)
	ListByUser(ctx context.Context, actorUserID uuid.UUID, limit, offset int) ([]*entity.AuditLog, error)
	ListByApplication(ctx context.Context, applicationID uuid.UUID, limit, offset int) ([]*entity.AuditLog, error)
}
