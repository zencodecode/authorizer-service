package repository

import (
	"context"
	"errors"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/auditlog"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/model"
	"gorm.io/gorm"
)

type auditLogRepository struct {
	db *gorm.DB
}

func NewAuditLogRepository(db *gorm.DB) auditlog.Repository {
	return &auditLogRepository{
		db: db,
	}
}

func (r *auditLogRepository) Create(ctx context.Context, log *entity.AuditLog) error {
	logModel := model.AuditLogFromEntity(log)
	return r.db.WithContext(ctx).Create(logModel).Error
}

func (r *auditLogRepository) GetByID(ctx context.Context, id string) (*entity.AuditLog, error) {
	var logModel model.AuditLog
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&logModel)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return logModel.ToEntity(), nil
}

func (r *auditLogRepository) ListByOrganization(ctx context.Context, organizationID string, limit, offset int) ([]*entity.AuditLog, error) {
	var logsModel []model.AuditLog
	err := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logsModel).Error
	if err != nil {
		return nil, err
	}

	logsEntity := make([]*entity.AuditLog, len(logsModel))
	for i, m := range logsModel {
		logsEntity[i] = m.ToEntity()
	}
	return logsEntity, nil
}

func (r *auditLogRepository) ListByUser(ctx context.Context, actorUserID string, limit, offset int) ([]*entity.AuditLog, error) {
	var logsModel []model.AuditLog
	err := r.db.WithContext(ctx).
		Where("actor_user_id = ?", actorUserID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logsModel).Error
	if err != nil {
		return nil, err
	}

	logsEntity := make([]*entity.AuditLog, len(logsModel))
	for i, m := range logsModel {
		logsEntity[i] = m.ToEntity()
	}
	return logsEntity, nil
}

func (r *auditLogRepository) ListByApplication(ctx context.Context, applicationID string, limit, offset int) ([]*entity.AuditLog, error) {
	var logsModel []model.AuditLog
	err := r.db.WithContext(ctx).
		Where("application_id = ?", applicationID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&logsModel).Error
	if err != nil {
		return nil, err
	}

	logsEntity := make([]*entity.AuditLog, len(logsModel))
	for i, m := range logsModel {
		logsEntity[i] = m.ToEntity()
	}
	return logsEntity, nil
}
