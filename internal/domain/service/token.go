package service

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

// JWTService handles JWT access token signing and validation.
// Refresh token persistence is handled by the oauthrefreshtoken repository.
type JWTService interface {
	GenerateAccessToken(ctx context.Context, claims *entity.Claims) (string, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*entity.Claims, error)
	GenerateRefreshToken() (string, error)
}
