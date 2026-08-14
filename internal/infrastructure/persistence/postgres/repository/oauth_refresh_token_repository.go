package repository

import (
	"context"
	"errors"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthrefreshtoken"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type oauthRefreshTokenRepository struct {
	db *gorm.DB
}

func NewOAuthRefreshTokenRepository(db *gorm.DB) oauthrefreshtoken.Repository {
	return &oauthRefreshTokenRepository{
		db: db,
	}
}

func (r *oauthRefreshTokenRepository) Create(ctx context.Context, token *entity.OAuthRefreshToken) error {
	tokenModel := model.OAuthRefreshTokenFromEntity(token)
	return r.db.WithContext(ctx).Create(tokenModel).Error
}

func (r *oauthRefreshTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.OAuthRefreshToken, error) {
	var tokenModel model.OAuthRefreshToken
	result := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&tokenModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, oauthrefreshtoken.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return tokenModel.ToEntity(), nil
}

func (r *oauthRefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.OAuthRefreshToken{}).
		Where("id = ?", id).
		Update("revoked_at", now).Error
}

func (r *oauthRefreshTokenRepository) RevokeByAccessTokenID(ctx context.Context, accessTokenID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.OAuthRefreshToken{}).
		Where("access_token_id = ? AND revoked_at IS NULL", accessTokenID).
		Update("revoked_at", now).Error
}

func (r *oauthRefreshTokenRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", before).
		Delete(&model.OAuthRefreshToken{}).Error
}
