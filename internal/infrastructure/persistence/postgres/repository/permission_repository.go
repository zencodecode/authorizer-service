package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/permission"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) permission.Repository {
	return &permissionRepository{
		db: db,
	}
}

func (r *permissionRepository) Create(ctx context.Context, perm *entity.Permission) error {
	permModel := model.PermissionFromEntity(perm)
	return r.db.WithContext(ctx).Create(permModel).Error
}

func (r *permissionRepository) GetByID(ctx context.Context, id string) (*entity.Permission, error) {
	var permModel model.Permission
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&permModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, permission.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return permModel.ToEntity(), nil
}

func (r *permissionRepository) GetBySlug(ctx context.Context, applicationID, slug string) (*entity.Permission, error) {
	var permModel model.Permission
	result := r.db.WithContext(ctx).
		Where("application_id = ? AND slug = ?", applicationID, slug).
		First(&permModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, permission.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return permModel.ToEntity(), nil
}

func (r *permissionRepository) Update(ctx context.Context, perm *entity.Permission) error {
	permModel := model.PermissionFromEntity(perm)
	return r.db.WithContext(ctx).Save(permModel).Error
}

func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Permission{}).Error
}

func (r *permissionRepository) ListByApplication(ctx context.Context, applicationID string, limit, offset int) ([]*entity.Permission, error) {
	var permsModel []model.Permission
	result := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Limit(limit).
		Offset(offset).
		Find(&permsModel)
	if result.Error != nil {
		return nil, result.Error
	}

	permsEntity := make([]*entity.Permission, len(permsModel))
	for i, m := range permsModel {
		permsEntity[i] = m.ToEntity()
	}
	return permsEntity, nil
}
