package oauthauthorizationcode

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, code *entity.OAuthAuthorizationCode) error
	GetByCodeHash(ctx context.Context, codeHash string) (*entity.OAuthAuthorizationCode, error)
	MarkUsed(ctx context.Context, id uuid.UUID) error
	DeleteExpired(ctx context.Context, before time.Time) error
}
