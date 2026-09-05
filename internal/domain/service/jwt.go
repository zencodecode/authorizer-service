package service

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type JWTService interface {
	GenerateAccessToken(ctx context.Context, claims *entity.Claims) (string, error)
	GenerateIDToken(ctx context.Context, claims *entity.IDTokenClaims) (string, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*entity.Claims, error)
	GenerateRefreshToken() (string, error)
}
