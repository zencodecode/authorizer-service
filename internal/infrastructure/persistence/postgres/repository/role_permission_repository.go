package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/rolepermission"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type rolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) rolepermission.Repository {
	return &rolePermissionRepository{
		db: db,
	}
}

func (r *rolePermissionRepository) AssignPermission(ctx context.Context, roleID, permissionID string) error {
	rolePermModel := model.RolePermission{
		RoleID:       parseUUID(roleID),
		PermissionID: parseUUID(permissionID),
	}
	return r.db.WithContext(ctx).Create(&rolePermModel).Error
}

func (r *rolePermissionRepository) RevokePermission(ctx context.Context, roleID, permissionID string) error {
	return r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&model.RolePermission{}).Error
}

func (r *rolePermissionRepository) HasPermission(ctx context.Context, roleID, permissionID string) (bool, error) {
	var rolePermModel model.RolePermission
	result := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		First(&rolePermModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (r *rolePermissionRepository) ListPermissionsByRole(ctx context.Context, roleID string) ([]*entity.Permission, error) {
	var permsModel []model.Permission
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permsModel).Error
	if err != nil {
		return nil, err
	}

	permsEntity := make([]*entity.Permission, len(permsModel))
	for i, m := range permsModel {
		permsEntity[i] = m.ToEntity()
	}
	return permsEntity, nil
}

func (r *rolePermissionRepository) ListRolesByPermission(ctx context.Context, permissionID string) ([]*entity.Role, error) {
	var rolesModel []model.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Where("role_permissions.permission_id = ?", permissionID).
		Find(&rolesModel).Error
	if err != nil {
		return nil, err
	}

	rolesEntity := make([]*entity.Role, len(rolesModel))
	for i, m := range rolesModel {
		rolesEntity[i] = m.ToEntity()
	}
	return rolesEntity, nil
}
