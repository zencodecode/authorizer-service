package applicationscope

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, scope *entity.ApplicationScope) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.ApplicationScope, error)
	GetByApplicationAndScope(ctx context.Context, applicationID uuid.UUID, scope string) (*entity.ApplicationScope, error)
	Update(ctx context.Context, scope *entity.ApplicationScope) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListByApplication(ctx context.Context, applicationID uuid.UUID) ([]*entity.ApplicationScope, error)
}
