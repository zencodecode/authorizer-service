package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/invitation"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) invitation.Repository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) Create(ctx context.Context, inv *entity.Invitation) error {
	m := model.InvitationFromEntity(inv)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *invitationRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Invitation, error) {
	var m model.Invitation
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, invitation.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *invitationRepository) GetByToken(ctx context.Context, token string) (*entity.Invitation, error) {
	var m model.Invitation
	result := r.db.WithContext(ctx).Where("token = ?", token).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, invitation.ErrNotFound
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *invitationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	result := r.db.WithContext(ctx).
		Model(&model.Invitation{}).
		Where("id = ?", id).
		Update("status", status)
	if result.RowsAffected == 0 {
		return invitation.ErrNotFound
	}
	return result.Error
}

func (r *invitationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Invitation{}).Error
}

func (r *invitationRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID, status *string, limit, offset int) ([]*entity.Invitation, error) {
	var models []model.Invitation
	query := r.db.WithContext(ctx).Where("organization_id = ?", organizationID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.Invitation, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
