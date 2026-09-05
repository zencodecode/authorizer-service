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
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthauthorizationcode"
)

type oauthAuthorizationCodeRepository struct {
	client *redis.Client
}

func NewOAuthAuthorizationCodeRepository(client *redis.Client) *oauthAuthorizationCodeRepository {
	return &oauthAuthorizationCodeRepository{client: client}
}

func (r *oauthAuthorizationCodeRepository) codeHashKey(codeHash string) string {
	return fmt.Sprintf("auth_code:hash:%s", codeHash)
}

func (r *oauthAuthorizationCodeRepository) idKey(id uuid.UUID) string {
	return fmt.Sprintf("auth_code:id:%s", id.String())
}

func (r *oauthAuthorizationCodeRepository) Create(ctx context.Context, code *entity.OAuthAuthorizationCode) error {
	data, err := json.Marshal(code)
	if err != nil {
		return fmt.Errorf("marshal authorization code: %w", err)
	}
	ttl := time.Until(code.ExpiresAt)
	if ttl <= 0 {
		return errors.New("authorization code already expired")
	}
	pipe := r.client.Pipeline()
	pipe.Set(ctx, r.codeHashKey(code.CodeHash), data, ttl)
	pipe.Set(ctx, r.idKey(code.ID), code.CodeHash, ttl)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *oauthAuthorizationCodeRepository) GetByCodeHash(ctx context.Context, codeHash string) (*entity.OAuthAuthorizationCode, error) {
	data, err := r.client.Get(ctx, r.codeHashKey(codeHash)).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("authorization code not found: %w", oauthauthorizationcode.ErrExpiredOrNotFound)
		}
		return nil, fmt.Errorf("get authorization code: %w", err)
	}
	var code entity.OAuthAuthorizationCode
	if err := json.Unmarshal(data, &code); err != nil {
		return nil, fmt.Errorf("unmarshal authorization code: %w", err)
	}
	return &code, nil
}

func (r *oauthAuthorizationCodeRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	codeHash, err := r.client.Get(ctx, r.idKey(id)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return fmt.Errorf("get code hash by id: %w", err)
	}
	pipe := r.client.Pipeline()
	pipe.Del(ctx, r.codeHashKey(codeHash))
	pipe.Del(ctx, r.idKey(id))
	_, err = pipe.Exec(ctx)
	return err
}

func (r *oauthAuthorizationCodeRepository) DeleteExpired(_ context.Context, _ time.Time) error {
	return nil
}
