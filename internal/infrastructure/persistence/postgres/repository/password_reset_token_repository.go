package repository

import (
	"context"
	"errors"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/passwordresettoken"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) passwordresettoken.Repository {
	return &passwordResetTokenRepository{
		db: db,
	}
}

func (r *passwordResetTokenRepository) Create(ctx context.Context, token *entity.PasswordResetToken) error {
	tokenModel := model.PasswordResetTokenFromEntity(token)
	return r.db.WithContext(ctx).Create(tokenModel).Error
}

func (r *passwordResetTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.PasswordResetToken, error) {
	var tokenModel model.PasswordResetToken
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

func (r *passwordResetTokenRepository) MarkUsed(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.PasswordResetToken{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}

func (r *passwordResetTokenRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", before).
		Delete(&model.PasswordResetToken{}).Error
}
