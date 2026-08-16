package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/rs256"
)

type jwtClaims struct {
	jwt.RegisteredClaims
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	OrgID   *string  `json:"org_id,omitempty"`
	OrgSlug *string  `json:"org_slug,omitempty"`
	Scopes  []string `json:"scopes"`

	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

type token struct {
	privateKey *rsa.PrivateKey
	expiry     time.Duration
	kid        string
	redis      *redis.Client
}

func NewTokenService(privateKey *rsa.PrivateKey, expiry time.Duration, redis *redis.Client) service.Token {
	return &token{privateKey: privateKey, expiry: expiry, redis: redis}
}

func (s *token) GenerateAccessToken(ctx context.Context, claims *entity.Claims) (string, error) {
	return rs256.GenerateAccessToken(s.privateKey, &jwtClaims{
		Name:        claims.Subject,
		Email:       claims.Email,
		OrgID:       claims.OrgID,
		OrgSlug:     claims.OrgSlug,
		Scopes:      claims.Scopes,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    claims.Issuer,
			Subject:   claims.Subject,
			Audience:  claims.Audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        claims.ID,
		},
	}, s.kid)
}

func (s *token) ValidateAccessToken(ctx context.Context, tokenString string) (*entity.Claims, error) {
	claims, err := rs256.ValidateAccessToken(tokenString, &s.privateKey.PublicKey, &jwtClaims{})
	if err != nil {
		return nil, err
	}

	return &entity.Claims{
		Name:        claims.Subject,
		Email:       claims.Email,
		OrgID:       claims.OrgID,
		OrgSlug:     claims.OrgSlug,
		Scopes:      claims.Scopes,
		Roles:       claims.Roles,
		Permissions: claims.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    claims.Issuer,
			Subject:   claims.Subject,
			Audience:  claims.Audience,
			ExpiresAt: claims.ExpiresAt,
			IssuedAt:  claims.IssuedAt,
			ID:        claims.ID,
		},
	}, nil
}

func (s *token) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32) // 256-bit
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (s *token) StoreRefreshToken(ctx context.Context, userID, token string) error {
	return s.redis.Set(ctx, "refresh:"+userID, token, 7*24*time.Hour).Err()
}

func (s *token) GetRefreshToken(ctx context.Context, userID string) (string, error) {
	return s.redis.Get(ctx, "refresh:"+userID).Result()
}

func (s *token) MapRefreshToUser(ctx context.Context, refreshToken, userID string) error {
	return s.redis.Set(ctx, "rt:"+refreshToken, userID, 7*24*time.Hour).Err()
}

func (s *token) GetUserIDByRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	return s.redis.Get(ctx, "rt:"+refreshToken).Result()
}

func (s *token) DeleteRefreshToken(ctx context.Context, userID string) error {
	return s.redis.Del(ctx, "refresh:"+userID).Err()
}
