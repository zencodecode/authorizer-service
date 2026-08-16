package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/auditlog"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) auditlog.Repository {
	return &auditLogRepository{db: db}
}

func (r *auditLogRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	m := model.AuditLogFromEntity(log)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *auditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.AuditLog, error) {
	var m model.AuditLog
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return m.ToEntity(), nil
}

func (r *auditLogRepository) ListByOrganization(ctx context.Context, organizationID uuid.UUID, limit, offset int) ([]*entity.AuditLog, error) {
	var models []model.AuditLog
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.AuditLog, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *auditLogRepository) ListByUser(ctx context.Context, actorUserID uuid.UUID, limit, offset int) ([]*entity.AuditLog, error) {
	var models []model.AuditLog
	err := r.db.WithContext(ctx).
		Where("actor_user_id = ?", actorUserID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.AuditLog, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}

func (r *auditLogRepository) ListByApplication(ctx context.Context, applicationID uuid.UUID, limit, offset int) ([]*entity.AuditLog, error) {
	var models []model.AuditLog
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.AuditLog, len(models))
	for i, m := range models {
		entities[i] = m.ToEntity()
	}
	return entities, nil
}
