package authorizesession

import (
	"context"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Save(ctx context.Context, challengeID string, session entity.AuthorizeSession, ttl time.Duration) error
	Get(ctx context.Context, challengeID string) (*entity.AuthorizeSession, error)
	Delete(ctx context.Context, challengeID string) error
}
