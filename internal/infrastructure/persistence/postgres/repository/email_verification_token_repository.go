package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/emailverificationtoken"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type emailVerificationTokenRepository struct {
	db *gorm.DB
}

func NewEmailVerificationTokenRepository(db *gorm.DB) emailverificationtoken.Repository {
	return &emailVerificationTokenRepository{db: db}
}

func (r *emailVerificationTokenRepository) Create(ctx context.Context, token *entity.EmailVerificationToken) error {
	m := model.EmailVerificationTokenFromEntity(token)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *emailVerificationTokenRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*entity.EmailVerificationToken, error) {
	var m model.EmailVerificationToken
	result := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, emailverificationtoken.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *emailVerificationTokenRepository) MarkUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&model.EmailVerificationToken{}).
		Where("id = ?", id).
		Update("used_at", now).Error
}

func (r *emailVerificationTokenRepository) DeleteExpired(ctx context.Context, before time.Time) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", before).
		Delete(&model.EmailVerificationToken{}).Error
}
