package applicationscope

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, scope *entity.ApplicationScope) error
	GetByID(ctx context.Context, id string) (*entity.ApplicationScope, error)
	GetByApplicationAndScope(ctx context.Context, applicationID, scope string) (*entity.ApplicationScope, error)
	Update(ctx context.Context, scope *entity.ApplicationScope) error
	Delete(ctx context.Context, id string) error
	ListByApplication(ctx context.Context, applicationID string) ([]*entity.ApplicationScope, error)
}
