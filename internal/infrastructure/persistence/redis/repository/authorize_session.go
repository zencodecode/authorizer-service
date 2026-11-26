package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/authorizesession"
)

type authorizeSessionRepository struct {
	client *redis.Client
}

func NewAuthorizeSessionRepository(client *redis.Client) authorizesession.Repository {
	return &authorizeSessionRepository{client: client}
}

func (r *authorizeSessionRepository) key(challengeID string) string {
	return fmt.Sprintf("authorize_session:%s", challengeID)
}

func (r *authorizeSessionRepository) Save(ctx context.Context, challengeID string, sess entity.AuthorizeSession, ttl time.Duration) error {
	data, err := json.Marshal(sess)
	if err != nil {
		return fmt.Errorf("marshal authorize session: %w", err)
	}
	return r.client.Set(ctx, r.key(challengeID), data, ttl).Err()
}

func (r *authorizeSessionRepository) Get(ctx context.Context, challengeID string) (*entity.AuthorizeSession, error) {
	data, err := r.client.Get(ctx, r.key(challengeID)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, authorizesession.ErrNotFound
		}
		return nil, fmt.Errorf("get authorize session: %w", err)
	}

	var sess entity.AuthorizeSession
	if err := json.Unmarshal(data, &sess); err != nil {
		return nil, fmt.Errorf("unmarshal authorize session: %w", err)
	}
	return &sess, nil
}

func (r *authorizeSessionRepository) Delete(ctx context.Context, challengeID string) error {
	return r.client.Del(ctx, r.key(challengeID)).Err()
}
