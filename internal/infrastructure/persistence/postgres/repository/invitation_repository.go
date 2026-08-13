package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/invitation"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) invitation.Repository {
	return &invitationRepository{
		db: db,
	}
}

func (r *invitationRepository) Create(ctx context.Context, inv *entity.Invitation) error {
	invModel := model.InvitationFromEntity(inv)
	return r.db.WithContext(ctx).Create(invModel).Error
}

func (r *invitationRepository) GetByID(ctx context.Context, id string) (*entity.Invitation, error) {
	var invModel model.Invitation
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&invModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return invModel.ToEntity(), nil
}

func (r *invitationRepository) GetByToken(ctx context.Context, token string) (*entity.Invitation, error) {
	var invModel model.Invitation
	result := r.db.WithContext(ctx).Where("token = ?", token).First(&invModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return invModel.ToEntity(), nil
}

func (r *invitationRepository) UpdateStatus(ctx context.Context, id, status string) error {
	result := r.db.WithContext(ctx).
		Model(&model.Invitation{}).
		Where("id = ?", id).
		Update("status", status)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r *invitationRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Invitation{}).Error
}

func (r *invitationRepository) ListByOrganization(ctx context.Context, organizationID string, status *string, limit, offset int) ([]*entity.Invitation, error) {
	var invsModel []model.Invitation
	query := r.db.WithContext(ctx).Where("organization_id = ?", organizationID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&invsModel).Error
	if err != nil {
		return nil, err
	}

	invsEntity := make([]*entity.Invitation, len(invsModel))
	for i, m := range invsModel {
		invsEntity[i] = m.ToEntity()
	}
	return invsEntity, nil
}
