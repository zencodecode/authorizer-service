package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type userRoleRepository struct {
	db *gorm.DB
}

func NewUserRoleRepository(db *gorm.DB) userrole.Repository {
	return &userRoleRepository{db: db}
}

func (r *userRoleRepository) AssignRole(ctx context.Context, userRole *entity.UserRole) error {
	m := model.UserRoleFromEntity(userRole)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *userRoleRepository) RevokeRole(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, roleID uuid.UUID) error {
	query := r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID)
	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	} else {
		query = query.Where("organization_id IS NULL")
	}
	return query.Delete(&model.UserRole{}).Error
}

func (r *userRoleRepository) HasRole(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, roleID uuid.UUID) (bool, error) {
	var m model.UserRole
	query := r.db.WithContext(ctx).Where("user_id = ? AND role_id = ?", userID, roleID)
	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	result := query.First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, userrole.ErrNotFound
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

func (r *userRoleRepository) ListRolesByUser(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID) ([]*entity.Role, error) {
	var models []model.Role
	query := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID)

	if organizationID != nil {
		query = query.Where("user_roles.organization_id = ?", *organizationID)
	}

	err := query.Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.Role, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *userRoleRepository) ListUsersByRole(ctx context.Context, roleID uuid.UUID, limit, offset int) ([]*entity.User, error) {
	var models []model.User
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Limit(limit).Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.User, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *userRoleRepository) ListPermissionsByUser(ctx context.Context, userID uuid.UUID, applicationID uuid.UUID, organizationID *uuid.UUID) ([]*entity.Permission, error) {
	var models []model.Permission
	query := r.db.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.application_id = ?", userID, applicationID)

	if organizationID != nil {
		query = query.Where("user_roles.organization_id = ?", *organizationID)
	}

	err := query.Distinct().Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.Permission, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
