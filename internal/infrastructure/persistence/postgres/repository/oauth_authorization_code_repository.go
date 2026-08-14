package repository

import (
	"context"
	"errors"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthauthorizationcode"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type oauthAuthorizationCodeRepository struct {
	db *gorm.DB
}

func NewOAuthAuthorizationCodeRepository(db *gorm.DB) oauthauthorizationcode.Repository {
	return &oauthAuthorizationCodeRepository{
		db: db,
	}
}

func (r *oauthAuthorizationCodeRepository) Create(ctx context.Context, code *entity.OAuthAuthorizationCode) error {
	codeModel := model.OAuthAuthorizationCodeFromEntity(code)
	return r.db.WithContext(ctx).Create(codeModel).Error
}

func (r *oauthAuthorizationCodeRepository) GetByCodeHash(ctx context.Context, codeHash string) (*entity.OAuthAuthorizationCode, error) {
	var codeModel model.OAuthAuthorizationCode
	result := r.db.WithContext(ctx).
		Where("code_hash = ?", codeHash).
		First(&codeModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, oauthauthorizationcode.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return codeModel.ToEntity(), nil
}

func (r *oauthAuthorizationCodeRepository) MarkUsed(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.OAuthAuthorizationCode{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}

func (r *oauthAuthorizationCodeRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", before).
		Delete(&model.OAuthAuthorizationCode{}).Error
}
