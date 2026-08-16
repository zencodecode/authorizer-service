package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/rolepermission"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type rolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) rolepermission.Repository {
	return &rolePermissionRepository{db: db}
}

func (r *rolePermissionRepository) AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	rp := model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}
	return r.db.WithContext(ctx).Create(&rp).Error
}

func (r *rolePermissionRepository) RevokePermission(ctx context.Context, roleID, permissionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&model.RolePermission{}).Error
}

func (r *rolePermissionRepository) HasPermission(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error) {
	var rp model.RolePermission
	result := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		First(&rp)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (r *rolePermissionRepository) ListPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*entity.Permission, error) {
	var models []model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.Permission, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *rolePermissionRepository) ListRolesByPermission(ctx context.Context, permissionID uuid.UUID) ([]*entity.Role, error) {
	var models []model.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Where("role_permissions.permission_id = ?", permissionID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.Role, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
