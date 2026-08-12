package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type AuditLog struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey"`
	OrganizationID *uuid.UUID      `gorm:"type:uuid"`
	ApplicationID  *uuid.UUID      `gorm:"type:uuid"`
	ActorUserID    *uuid.UUID      `gorm:"type:uuid"`
	Action         string          `gorm:"type:varchar(255);not null"`
	ResourceType   *string         `gorm:"type:varchar(255)"`
	ResourceID     *uuid.UUID      `gorm:"type:uuid"`
	Metadata       json.RawMessage `gorm:"type:jsonb"`
	CreatedAt      time.Time       `gorm:"not null;default:now()"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func AuditLogFromEntity(e *entity.AuditLog) *AuditLog {
	var metadata json.RawMessage
	if e.Metadata != nil {
		metadata, _ = json.Marshal(e.Metadata)
	}

	return &AuditLog{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		ApplicationID:  e.ApplicationID,
		ActorUserID:    e.ActorUserID,
		Action:         e.Action,
		ResourceType:   e.ResourceType,
		ResourceID:     e.ResourceID,
		Metadata:       metadata,
		CreatedAt:      e.CreatedAt,
	}
}

func (m *AuditLog) ToEntity() *entity.AuditLog {
	var metadata map[string]any
	if m.Metadata != nil {
		_ = json.Unmarshal(m.Metadata, &metadata)
	}

	return &entity.AuditLog{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		ApplicationID:  m.ApplicationID,
		ActorUserID:    m.ActorUserID,
		Action:         m.Action,
		ResourceType:   m.ResourceType,
		ResourceID:     m.ResourceID,
		Metadata:       metadata,
		CreatedAt:      m.CreatedAt,
	}
}
