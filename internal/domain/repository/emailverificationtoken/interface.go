package emailverificationtoken

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, token *entity.EmailVerificationToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.EmailVerificationToken, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context, before time.Time) error
}
