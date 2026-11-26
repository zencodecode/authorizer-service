package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/applicationscope"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type applicationScopeRepository struct {
	db *gorm.DB
}

func NewApplicationScopeRepository(db *gorm.DB) applicationscope.Repository {
	return &applicationScopeRepository{db: db}
}

func (r *applicationScopeRepository) Create(ctx context.Context, scope *entity.ApplicationScope) error {
	m := model.ApplicationScopeFromEntity(scope)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *applicationScopeRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.ApplicationScope, error) {
	var m model.ApplicationScope
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, applicationscope.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *applicationScopeRepository) GetByApplicationAndScope(ctx context.Context, applicationID uuid.UUID, scope string) (*entity.ApplicationScope, error) {
	var m model.ApplicationScope
	result := r.db.WithContext(ctx).
		Where("application_id = ? AND scope = ?", applicationID, scope).
		First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, applicationscope.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *applicationScopeRepository) Update(ctx context.Context, scope *entity.ApplicationScope) error {
	m := model.ApplicationScopeFromEntity(scope)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *applicationScopeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ApplicationScope{}).Error
}

func (r *applicationScopeRepository) ListByApplication(ctx context.Context, applicationID uuid.UUID) ([]*entity.ApplicationScope, error) {
	var models []model.ApplicationScope
	result := r.db.WithContext(ctx).Where("application_id = ?", applicationID).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*entity.ApplicationScope, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
