package oauthrefreshtoken

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, token *entity.OAuthRefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*entity.OAuthRefreshToken, error)
	Revoke(ctx context.Context, id uuid.UUID) error
	RevokeByAccessTokenID(ctx context.Context, accessTokenID uuid.UUID) error
	RevokeAllByUserAndApplication(ctx context.Context, userID, applicationID uuid.UUID) error
	RevokeAllByUser(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context, before time.Time) error
}
