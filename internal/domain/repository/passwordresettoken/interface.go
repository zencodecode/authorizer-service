package passwordresettoken

import (
	"context"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, token *entity.PasswordResetToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context, before time.Time) error
}
