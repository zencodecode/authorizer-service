package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *entity.User) error {
	m := model.UserFromEntity(u)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var m model.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, user.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, user.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *userRepository) Update(ctx context.Context, u *entity.User) error {
	m := model.UserFromEntity(u)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *userRepository) UpdateEmailVerified(ctx context.Context, id uuid.UUID, verifiedAt *time.Time) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("email_verified_at", verifiedAt).Error
}

func (r *userRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var models []model.User
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*entity.User, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
