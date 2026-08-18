package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/role"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) role.Repository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, rl *entity.Role) error {
	m := model.RoleFromEntity(rl)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Role, error) {
	var m model.Role
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, role.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *roleRepository) GetBySlug(ctx context.Context, organizationID *uuid.UUID, applicationID uuid.UUID, slug string) (*entity.Role, error) {
	var m model.Role
	query := r.db.WithContext(ctx).Where("application_id = ? AND slug = ?", applicationID, slug)
	if organizationID != nil {
		query = query.Where("organization_id = ?", *organizationID)
	} else {
		query = query.Where("organization_id IS NULL")
	}

	result := query.First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, role.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *roleRepository) Update(ctx context.Context, rl *entity.Role) error {
	m := model.RoleFromEntity(rl)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *roleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Role{}).Error
}

func (r *roleRepository) ListByApplication(ctx context.Context, applicationID uuid.UUID, limit, offset int) ([]*entity.Role, error) {
	var models []model.Role
	result := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Limit(limit).Offset(offset).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*entity.Role, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *roleRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entity.Role, error) {
	var models []model.Role
	result := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Limit(limit).Offset(offset).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*entity.Role, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *roleRepository) ListByOrganizationAndApplication(ctx context.Context, organizationID, applicationID uuid.UUID) ([]*entity.Role, error) {
	var models []model.Role
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND application_id = ?", organizationID, applicationID).
		Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*entity.Role, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
