package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/applicationscope"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type applicationScopeRepository struct {
	db *gorm.DB
}

func NewApplicationScopeRepository(db *gorm.DB) applicationscope.Repository {
	return &applicationScopeRepository{
		db: db,
	}
}

func (r *applicationScopeRepository) Create(ctx context.Context, scope *entity.ApplicationScope) error {
	scopeModel := model.ApplicationScopeFromEntity(scope)
	return r.db.WithContext(ctx).Create(scopeModel).Error
}

func (r *applicationScopeRepository) GetByID(ctx context.Context, id string) (*entity.ApplicationScope, error) {
	var scopeModel model.ApplicationScope
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&scopeModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return scopeModel.ToEntity(), nil
}

func (r *applicationScopeRepository) GetByApplicationAndScope(
	ctx context.Context,
	applicationID,
	scope string,
) (*entity.ApplicationScope, error) {
	var scopeModel model.ApplicationScope
	result := r.db.WithContext(ctx).
		Where("application_id = ? AND scope = ?", applicationID, scope).
		First(&scopeModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return scopeModel.ToEntity(), nil
}

func (r *applicationScopeRepository) Update(ctx context.Context, scope *entity.ApplicationScope) error {
	scopeModel := model.ApplicationScopeFromEntity(scope)
	return r.db.WithContext(ctx).Save(scopeModel).Error
}

func (r *applicationScopeRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.ApplicationScope{}).Error
}

func (r *applicationScopeRepository) ListByApplication(ctx context.Context, applicationID string) ([]*entity.ApplicationScope, error) {
	var scopesModel []model.ApplicationScope
	result := r.db.WithContext(ctx).Find(&scopesModel)
	if result.Error != nil {
		return nil, result.Error
	}

	scopesEntity := make([]*entity.ApplicationScope, len(scopesModel))
	for i, m := range scopesModel {
		scopesEntity[i] = m.ToEntity()
	}
	return scopesEntity, nil
}
