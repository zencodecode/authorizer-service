package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) userrole.Repository {
	return &userRoleRepository{
		db: db,
	}
}

func (r *userRoleRepository) AssignRole(ctx context.Context, userRole *entity.UserRole) error {
	userRoleModel := model.UserRoleFromEntity(userRole)
	return r.db.WithContext(ctx).Create(userRoleModel).Error
}

func (r *userRoleRepository) RevokeRole(ctx context.Context, userID, organizationID, roleID string) error {
	query := r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID)
	if organizationID != "" {
		query = query.Where("organization_id = ?", organizationID)
	} else {
		query = query.Where("organization_id IS NULL")
	}
	return query.Delete(&model.UserRole{}).Error
}

func (r *userRoleRepository) HasRole(ctx context.Context, userID, organizationID, roleID string) (bool, error) {
	var userRoleModel model.UserRole
	query := r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID)
	if organizationID != "" {
		query = query.Where("organization_id = ?", organizationID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	result := query.First(&userRoleModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (r *userRoleRepository) ListRolesByUser(ctx context.Context, userID string, organizationID *string) ([]*entity.Role, error) {
	var rolesModel []model.Role
	query := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID)

	if organizationID != nil {
		query = query.Where("user_roles.organization_id = ?", *organizationID)
	}

	err := query.Find(&rolesModel).Error
	if err != nil {
		return nil, err
	}

	rolesEntity := make([]*entity.Role, len(rolesModel))
	for i, m := range rolesModel {
		rolesEntity[i] = m.ToEntity()
	}
	return rolesEntity, nil
}

func (r *userRoleRepository) ListUsersByRole(ctx context.Context, roleID string, limit, offset int) ([]*entity.User, error) {
	var usersModel []model.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Limit(limit).
		Offset(offset).
		Find(&usersModel).Error
	if err != nil {
		return nil, err
	}

	usersEntity := make([]*entity.User, len(usersModel))
	for i, m := range usersModel {
		usersEntity[i] = m.ToEntity()
	}
	return usersEntity, nil
}

func (r *userRoleRepository) ListPermissionsByUser(ctx context.Context, userID, applicationID string, organizationID *string) ([]*entity.Permission, error) {
	var permsModel []model.Permission
	query := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.application_id = ?", userID, applicationID)

	if organizationID != nil {
		query = query.Where("user_roles.organization_id = ?", *organizationID)
	}

	err := query.Distinct().Find(&permsModel).Error
	if err != nil {
		return nil, err
	}

	permsEntity := make([]*entity.Permission, len(permsModel))
	for i, m := range permsModel {
		permsEntity[i] = m.ToEntity()
	}
	return permsEntity, nil
}
