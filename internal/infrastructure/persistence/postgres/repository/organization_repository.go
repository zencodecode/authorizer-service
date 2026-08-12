package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organization"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) organization.Repository {
	return &organizationRepository{
		db: db,
	}
}

func (r *organizationRepository) Create(ctx context.Context, org *entity.Organization) error {
	orgModel := model.OrganizationFromEntity(org)
	return r.db.WithContext(ctx).Create(orgModel).Error
}

func (r *organizationRepository) GetByID(ctx context.Context, id string) (*entity.Organization, error) {
	var orgModel model.Organization
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&orgModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return orgModel.ToEntity(), nil
}

func (r *organizationRepository) GetBySlug(ctx context.Context, slug string) (*entity.Organization, error) {
	var orgModel model.Organization
	result := r.db.WithContext(ctx).Where("slug = ?", slug).First(&orgModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return orgModel.ToEntity(), nil
}

func (r *organizationRepository) Update(ctx context.Context, org *entity.Organization) error {
	orgModel := model.OrganizationFromEntity(org)
	return r.db.WithContext(ctx).Save(orgModel).Error
}

func (r *organizationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Organization{}).Error
}

func (r *organizationRepository) List(ctx context.Context, limit, offset int) ([]*entity.Organization, error) {
	var orgsModel []model.Organization
	result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&orgsModel)
	if result.Error != nil {
		return nil, result.Error
	}

	orgsEntity := make([]*entity.Organization, len(orgsModel))
	for i, m := range orgsModel {
		orgsEntity[i] = m.ToEntity()
	}
	return orgsEntity, nil
}
