package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/permission"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) permission.Repository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(ctx context.Context, perm *entity.Permission) error {
	m := model.PermissionFromEntity(perm)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *permissionRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Permission, error) {
	var m model.Permission
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, permission.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *permissionRepository) GetBySlug(ctx context.Context, applicationID uuid.UUID, slug string) (*entity.Permission, error) {
	var m model.Permission
	result := r.db.WithContext(ctx).
		Where("application_id = ? AND slug = ?", applicationID, slug).
		First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, permission.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *permissionRepository) Update(ctx context.Context, perm *entity.Permission) error {
	m := model.PermissionFromEntity(perm)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Permission{}).Error
}

func (r *permissionRepository) ListByApplication(ctx context.Context, applicationID uuid.UUID, limit, offset int) ([]*entity.Permission, error) {
	var models []model.Permission
	result := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Limit(limit).Offset(offset).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*entity.Permission, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
