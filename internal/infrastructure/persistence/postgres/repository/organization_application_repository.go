package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationapplication"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type organizationApplicationRepository struct {
	db *gorm.DB
}

func NewOrganizationApplicationRepository(db *gorm.DB) organizationapplication.Repository {
	return &organizationApplicationRepository{
		db: db,
	}
}

func (r *organizationApplicationRepository) Create(ctx context.Context, orgApp *entity.OrganizationApplication) error {
	orgAppModel := model.OrganizationApplicationFromEntity(orgApp)
	return r.db.WithContext(ctx).Create(orgAppModel).Error
}

func (r *organizationApplicationRepository) GetByOrganizationAndApplication(
	ctx context.Context,
	organizationID,
	applicationID string,
) (*entity.OrganizationApplication, error) {
	var orgAppModel model.OrganizationApplication
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND application_id = ?", organizationID, applicationID).
		First(&orgAppModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, organizationapplication.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return orgAppModel.ToEntity(), nil
}

func (r *organizationApplicationRepository) SetActive(ctx context.Context, organizationID, applicationID string, isActive bool) error {
	result := r.db.WithContext(ctx).
		Model(&model.OrganizationApplication{}).
		Where("organization_id = ? AND application_id = ?", organizationID, applicationID).
		Update("is_active", isActive)

	if result.RowsAffected == 0 {
		return organizationapplication.ErrNotFound
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *organizationApplicationRepository) Delete(ctx context.Context, organizationID, applicationID string) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND application_id = ?", organizationID, applicationID).
		Delete(&model.OrganizationApplication{}).Error
}

func (r *organizationApplicationRepository) ListApplicationsByOrganization(
	ctx context.Context,
	organizationID string,
) ([]*entity.OrganizationApplication, error) {
	var orgAppsModel []model.OrganizationApplication
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Find(&orgAppsModel).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.OrganizationApplication, len(orgAppsModel))
	for i, m := range orgAppsModel {
		result[i] = m.ToEntity()
	}
	return result, nil
}

func (r *organizationApplicationRepository) ListOrganizationsByApplication(
	ctx context.Context,
	applicationID string,
	limit,
	offset int,
) ([]*entity.OrganizationApplication, error) {
	var orgAppsModel []model.OrganizationApplication
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Find(&orgAppsModel).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.OrganizationApplication, len(orgAppsModel))
	for i, m := range orgAppsModel {
		result[i] = m.ToEntity()
	}
	return result, nil
}
