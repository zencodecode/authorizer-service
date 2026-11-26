package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) application.Repository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, app *entity.Application) error {
	m := model.ApplicationFromEntity(app)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *applicationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Application, error) {
	var m model.Application
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, application.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *applicationRepository) GetBySlug(ctx context.Context, slug string) (*entity.Application, error) {
	var m model.Application
	result := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, application.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *applicationRepository) GetByClientID(ctx context.Context, clientID string) (*entity.Application, error) {
	var m model.Application
	result := r.db.WithContext(ctx).Where("client_id = ?", clientID).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, application.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *applicationRepository) Update(ctx context.Context, app *entity.Application) error {
	m := model.ApplicationFromEntity(app)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *applicationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Application{}).Error
}

func (r *applicationRepository) List(ctx context.Context, limit, offset int) ([]*entity.Application, error) {
	var models []model.Application
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*entity.Application, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
