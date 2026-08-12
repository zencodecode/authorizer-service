package oauthaccesstoken

import (
	"context"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, token *entity.OAuthAccessToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.OAuthAccessToken, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllByUser(ctx context.Context, userID string) error
	RevokeAllByUserAndApplication(ctx context.Context, userID, applicationID string) error
	DeleteExpired(ctx context.Context, before time.Time) error
}
