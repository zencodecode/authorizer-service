package service

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Token interface {
	GenerateAccessToken(ctx context.Context, claims *entity.Claims) (string, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*entity.Claims, error)
	GenerateRefreshToken() (string, error)
	StoreRefreshToken(ctx context.Context, userID, token string) error
	GetRefreshToken(ctx context.Context, userID string) (string, error)
	MapRefreshToUser(ctx context.Context, refreshToken, userID string) error
	GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error)
	DeleteRefreshToken(ctx context.Context, userID string) error
}
