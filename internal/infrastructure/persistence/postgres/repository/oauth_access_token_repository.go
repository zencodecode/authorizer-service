package repository

import (
	"context"
	"errors"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthaccesstoken"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type oauthAccessTokenRepository struct {
	db *gorm.DB
}

func NewOAuthAccessTokenRepository(db *gorm.DB) oauthaccesstoken.Repository {
	return &oauthAccessTokenRepository{
		db: db,
	}
}

func (r *oauthAccessTokenRepository) Create(ctx context.Context, token *entity.OAuthAccessToken) error {
	tokenModel := model.OAuthAccessTokenFromEntity(token)
	return r.db.WithContext(ctx).Create(tokenModel).Error
}

func (r *oauthAccessTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.OAuthAccessToken, error) {
	var tokenModel model.OAuthAccessToken
	result := r.db.WithContext(ctx).
		Where("token_hash = ?", tokenHash).
		First(&tokenModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return tokenModel.ToEntity(), nil
}

func (r *oauthAccessTokenRepository) Revoke(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.OAuthAccessToken{}).
		Where("id = ?", id).
		Update("revoked_at", now).Error
}

func (r *oauthAccessTokenRepository) RevokeAllByUser(ctx context.Context, userID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.OAuthAccessToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

func (r *oauthAccessTokenRepository) RevokeAllByUserAndApplication(ctx context.Context, userID, applicationID string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.OAuthAccessToken{}).
		Where("user_id = ? AND application_id = ? AND revoked_at IS NULL", userID, applicationID).
		Update("revoked_at", now).Error
}

func (r *oauthAccessTokenRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", before).
		Delete(&model.OAuthAccessToken{}).Error
}
