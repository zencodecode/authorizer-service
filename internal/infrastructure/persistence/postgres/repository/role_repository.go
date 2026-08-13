package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/role"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) role.Repository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	roleModel := model.RoleFromEntity(role)
	return r.db.WithContext(ctx).Create(roleModel).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id string) (*entity.Role, error) {
	var roleModel model.Role
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&roleModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return roleModel.ToEntity(), nil
}

func (r *roleRepository) GetBySlug(ctx context.Context, organizationID *string, applicationID, slug string) (*entity.Role, error) {
	var roleModel model.Role
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND application_id = ? AND slug = ?", organizationID, applicationID, slug).
		First(&roleModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return roleModel.ToEntity(), nil
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	roleModel := model.RoleFromEntity(role)
	return r.db.WithContext(ctx).Save(roleModel).Error
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Role{}).Error
}

func (r *roleRepository) ListByApplication(ctx context.Context, applicationID string, limit, offset int) ([]*entity.Role, error) {
	var rolesModel []model.Role
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&rolesModel)
	if result.Error != nil {
		return nil, result.Error
	}

	rolesEntity := make([]*entity.Role, len(rolesModel))
	for i, m := range rolesModel {
		rolesEntity[i] = m.ToEntity()
	}
	return rolesEntity, nil
}

func (r *roleRepository) ListByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*entity.Role, error) {
	var rolesModel []model.Role
	result := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Limit(limit).
		Offset(offset).
		Find(&rolesModel)
	if result.Error != nil {
		return nil, result.Error
	}

	rolesEntity := make([]*entity.Role, len(rolesModel))
	for i, m := range rolesModel {
		rolesEntity[i] = m.ToEntity()
	}
	return rolesEntity, nil
}

func (r *roleRepository) ListByOrganizationAndApplication(ctx context.Context, organizationID, applicationID string) ([]*entity.Role, error) {
	var rolesModel []model.Role
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND application_id = ?", organizationID, applicationID).
		Find(&rolesModel)
	if result.Error != nil {
		return nil, result.Error
	}

	rolesEntity := make([]*entity.Role, len(rolesModel))
	for i, m := range rolesModel {
		rolesEntity[i] = m.ToEntity()
	}
	return rolesEntity, nil
}
