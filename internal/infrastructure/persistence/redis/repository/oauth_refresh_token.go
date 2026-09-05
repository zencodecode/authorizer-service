package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthrefreshtoken"
)

type oauthRefreshTokenRepository struct {
	client *redis.Client
}

func NewOAuthRefreshTokenRepository(client *redis.Client) *oauthRefreshTokenRepository {
	return &oauthRefreshTokenRepository{client: client}
}

func (r *oauthRefreshTokenRepository) tokenHashKey(tokenHash string) string {
	return fmt.Sprintf("refresh_token:hash:%s", tokenHash)
}

func (r *oauthRefreshTokenRepository) idKey(id uuid.UUID) string {
	return fmt.Sprintf("refresh_token:id:%s", id.String())
}

func (r *oauthRefreshTokenRepository) userAppKey(userID, appID uuid.UUID) string {
	return fmt.Sprintf("refresh_token:user:%s:app:%s", userID.String(), appID.String())
}

func (r *oauthRefreshTokenRepository) userKey(userID uuid.UUID) string {
	return fmt.Sprintf("refresh_token:user:%s", userID.String())
}

func (r *oauthRefreshTokenRepository) Create(ctx context.Context, token *entity.OAuthRefreshToken) error {
	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return errors.New("refresh token already expired")
	}
	data, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("marshal refresh token: %w", err)
	}
	pipe := r.client.Pipeline()
	pipe.Set(ctx, r.tokenHashKey(token.TokenHash), data, ttl)
	pipe.Set(ctx, r.idKey(token.ID), token.TokenHash, ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *oauthRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.OAuthRefreshToken, error) {
	data, err := r.client.Get(ctx, r.tokenHashKey(tokenHash)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, oauthrefreshtoken.ErrNotFound
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	var token entity.OAuthRefreshToken
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("unmarshal refresh token: %w", err)
	}
	return &token, nil
}

func (r *oauthRefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	tokenHash, err := r.client.Get(ctx, r.idKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("get token hash by id: %w", err)
	}
	pipe := r.client.Pipeline()
	pipe.Del(ctx, r.tokenHashKey(tokenHash))
	pipe.Del(ctx, r.idKey(id))
	_, err = pipe.Exec(ctx)
	return err
}

func (r *oauthRefreshTokenRepository) RevokeByAccessTokenID(ctx context.Context, accessTokenID uuid.UUID) error {
	return nil
}

func (r *oauthRefreshTokenRepository) RevokeAllByUserAndApplication(ctx context.Context, userID, applicationID uuid.UUID) error {
	setKey := r.userAppKey(userID, applicationID)
	members, err := r.client.SMembers(ctx, setKey).Result()
	if err != nil {
		return fmt.Errorf("get user refresh tokens: %w", err)
	}
	if len(members) == 0 {
		return nil
	}
	pipe := r.client.Pipeline()
	for _, tokenHash := range members {
		pipe.Del(ctx, r.tokenHashKey(tokenHash))
	}
	pipe.Del(ctx, setKey)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *oauthRefreshTokenRepository) RevokeAllByUser(ctx context.Context, userID uuid.UUID) error {
	setKey := r.userKey(userID)
	members, err := r.client.SMembers(ctx, setKey).Result()
	if err != nil {
		return fmt.Errorf("get user refresh tokens: %w", err)
	}
	if len(members) == 0 {
		return nil
	}
	pipe := r.client.Pipeline()
	for _, tokenHash := range members {
		pipe.Del(ctx, r.tokenHashKey(tokenHash))
	}
	pipe.Del(ctx, setKey)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *oauthRefreshTokenRepository) DeleteExpired(_ context.Context, _ time.Time) error {
	return nil
}
