package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type OrganizationUser struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_user"`
	UserID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_org_user"`
	Status         string    `gorm:"type:varchar(50);not null;default:'active'"`
	JoinedAt       time.Time `gorm:"type:timestamptz;not null;default:now()"`
	CreatedAt      time.Time `gorm:"not null;default:now()"`
}

func (OrganizationUser) TableName() string {
	return "organization_users"
}

func OrganizationUserFromEntity(e *entity.OrganizationUser) *OrganizationUser {
	return &OrganizationUser{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		UserID:         e.UserID,
		Status:         e.Status,
		JoinedAt:       e.JoinedAt,
		CreatedAt:      e.CreatedAt,
	}
}

func (m *OrganizationUser) ToEntity() *entity.OrganizationUser {
	return &entity.OrganizationUser{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		Status:         m.Status,
		JoinedAt:       m.JoinedAt,
		CreatedAt:      m.CreatedAt,
	}
}
