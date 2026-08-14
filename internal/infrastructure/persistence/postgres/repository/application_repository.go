package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) application.Repository {
	return &applicationRepository{
		db: db,
	}
}

func (r *applicationRepository) Create(ctx context.Context, app *entity.Application) error {
	appModel := model.ApplicationFromEntity(app)
	return r.db.WithContext(ctx).Create(appModel).Error
}

func (r *applicationRepository) GetByID(ctx context.Context, id string) (*entity.Application, error) {
	var appModel model.Application
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&appModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, application.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return appModel.ToEntity(), nil
}

func (r *applicationRepository) GetBySlug(ctx context.Context, slug string) (*entity.Application, error) {
	var appModel model.Application
	result := r.db.WithContext(ctx).Where("slug = ?", slug).First(&appModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, application.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return appModel.ToEntity(), nil
}

func (r *applicationRepository) GetByClientID(ctx context.Context, clientID string) (*entity.Application, error) {
	var appModel model.Application
	result := r.db.WithContext(ctx).Where("client_id = ?", clientID).First(&appModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, application.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return appModel.ToEntity(), nil
}

func (r *applicationRepository) Update(ctx context.Context, app *entity.Application) error {
	appModel := model.ApplicationFromEntity(app)
	return r.db.WithContext(ctx).Save(appModel).Error
}

func (r *applicationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Application{}).Error
}

func (r *applicationRepository) List(ctx context.Context, limit, offset int) ([]*entity.Application, error) {
	var appsModel []model.Application
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&appsModel)
	if result.Error != nil {
		return nil, result.Error
	}

	appsEntity := make([]*entity.Application, len(appsModel))
	for i, m := range appsModel {
		appsEntity[i] = m.ToEntity()
	}
	return appsEntity, nil
}
