package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationuser"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type organizationUserRepository struct {
	db *gorm.DB
}

func NewOrganizationUserRepository(db *gorm.DB) organizationuser.Repository {
	return &organizationUserRepository{
		db: db,
	}
}

func (r *organizationUserRepository) Create(ctx context.Context, orgUser *entity.OrganizationUser) error {
	orgUserModel := model.OrganizationUserFromEntity(orgUser)
	return r.db.WithContext(ctx).Create(orgUserModel).Error
}

func (r *organizationUserRepository) GetByOrganizationAndUser(ctx context.Context, organizationID, userID string) (
	*entity.OrganizationUser,
	error,
) {
	var orgUserModel model.OrganizationUser
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		First(&orgUserModel)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return orgUserModel.ToEntity(), nil
}

func (r *organizationUserRepository) UpdateStatus(
	ctx context.Context,
	organizationID,
	userID,
	status string,
) error {
	result := r.db.WithContext(ctx).
		Model(&entity.OrganizationUser{}).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Update("status", status)

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *organizationUserRepository) Delete(ctx context.Context, organizationID, userID string) error {
	return r.db.WithContext(ctx).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Delete(&model.OrganizationUser{}).Error
}

func (r *organizationUserRepository) ListUsersByOrganization(
	ctx context.Context,
	organizationID string,
	limit,
	offset int,
) ([]*entity.OrganizationUser, error) {
	var orgUsersModel []model.OrganizationUser
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Find(&orgUsersModel).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.OrganizationUser, len(orgUsersModel))
	for i, m := range orgUsersModel {
		result[i] = m.ToEntity()
	}
	return result, nil
}

func (r *organizationUserRepository) ListOrganizationsByUser(ctx context.Context, userID string) ([]*entity.OrganizationUser, error) {
	var orgUsersModel []model.OrganizationUser
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&orgUsersModel).Error
	if err != nil {
		return nil, err
	}

	result := make([]*entity.OrganizationUser, len(orgUsersModel))
	for i, m := range orgUsersModel {
		result[i] = m.ToEntity()
	}
	return result, nil
}
