package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationapplication"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type organizationApplicationRepository struct {
	db *gorm.DB
}

func NewOrganizationApplicationRepository(db *gorm.DB) organizationapplication.Repository {
	return &organizationApplicationRepository{db: db}
}

func (r *organizationApplicationRepository) Create(ctx context.Context, orgApp *entity.OrganizationApplication) error {
	m := model.OrganizationApplicationFromEntity(orgApp)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *organizationApplicationRepository) GetByOrganizationAndApplication(ctx context.Context, organizationID, applicationID uuid.UUID) (*entity.OrganizationApplication, error) {
	var m model.OrganizationApplication
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND application_id = ?", organizationID, applicationID).
		First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *organizationApplicationRepository) SetActive(ctx context.Context, organizationID, applicationID uuid.UUID, isActive bool) error {
	result := r.db.WithContext(ctx).
		Model(&model.OrganizationApplication{}).
		Where("organization_id = ? AND application_id = ?", organizationID, applicationID).
		Update("is_active", isActive)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *organizationApplicationRepository) Delete(ctx context.Context, organizationID, applicationID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND application_id = ?", organizationID, applicationID).
		Delete(&model.OrganizationApplication{}).Error
}

func (r *organizationApplicationRepository) ListApplicationsByOrganization(ctx context.Context, organizationID uuid.UUID) ([]*entity.OrganizationApplication, error) {
	var models []model.OrganizationApplication
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.OrganizationApplication, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *organizationApplicationRepository) ListOrganizationsByApplication(ctx context.Context, applicationID uuid.UUID, limit, offset int) ([]*entity.OrganizationApplication, error) {
	var models []model.OrganizationApplication
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Limit(limit).Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.OrganizationApplication, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
