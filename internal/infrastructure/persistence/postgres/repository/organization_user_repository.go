package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationuser"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type organizationUserRepository struct {
	db *gorm.DB
}

func NewOrganizationUserRepository(db *gorm.DB) organizationuser.Repository {
	return &organizationUserRepository{db: db}
}

func (r *organizationUserRepository) Create(ctx context.Context, orgUser *entity.OrganizationUser) error {
	m := model.OrganizationUserFromEntity(orgUser)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *organizationUserRepository) GetByOrganizationAndUser(ctx context.Context, organizationID, userID uuid.UUID) (*entity.OrganizationUser, error) {
	var m model.OrganizationUser
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, organizationuser.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *organizationUserRepository) UpdateStatus(ctx context.Context, organizationID, userID uuid.UUID, status string) error {
	result := r.db.WithContext(ctx).
		Model(&model.OrganizationUser{}).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Update("status", status)
	if result.RowsAffected == 0 {
		return organizationuser.ErrNotFound
	}
	return result.Error
}

func (r *organizationUserRepository) Delete(ctx context.Context, organizationID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Delete(&model.OrganizationUser{}).Error
}

func (r *organizationUserRepository) ListUsersByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entity.OrganizationUser, error) {
	var models []model.OrganizationUser
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Limit(limit).Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.OrganizationUser, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *organizationUserRepository) ListOrganizationsByUser(ctx context.Context, userID uuid.UUID) ([]*entity.OrganizationUser, error) {
	var models []model.OrganizationUser
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.OrganizationUser, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
