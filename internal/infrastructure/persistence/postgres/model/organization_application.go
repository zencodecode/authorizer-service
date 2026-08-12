package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type OrganizationApplication struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_app"`
	ApplicationID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_app"`
	IsActive       bool      `gorm:"not null;default:true"`
	ActivatedAt    time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

func (OrganizationApplication) TableName() string {
	return "organization_applications"
}

func OrganizationApplicationFromEntity(e *entity.OrganizationApplication) *OrganizationApplication {
	return &OrganizationApplication{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		ApplicationID:  e.ApplicationID,
		IsActive:       e.IsActive,
		ActivatedAt:    e.ActivatedAt,
	}
}

func (m *OrganizationApplication) ToEntity() *entity.OrganizationApplication {
	return &entity.OrganizationApplication{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		ApplicationID:  m.ApplicationID,
		IsActive:       m.IsActive,
		ActivatedAt:    m.ActivatedAt,
	}
}
