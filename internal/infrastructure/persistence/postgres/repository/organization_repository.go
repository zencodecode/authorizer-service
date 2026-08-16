package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organization"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) organization.Repository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(ctx context.Context, org *entity.Organization) error {
	m := model.OrganizationFromEntity(org)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Organization, error) {
	var m model.Organization
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*entity.Organization, error) {
	var m model.Organization
	result := r.db.WithContext(ctx).Where("slug = ?", slug).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *organizationRepository) Update(ctx context.Context, org *entity.Organization) error {
	m := model.OrganizationFromEntity(org)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Organization{}).Error
}

func (r *organizationRepository) List(ctx context.Context, limit, offset int) ([]*entity.Organization, error) {
	var models []model.Organization
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	entities := make([]*entity.Organization, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
