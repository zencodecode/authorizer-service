package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Role struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_org_app_role_slug"`
	ApplicationID  uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_org_app_role_slug"`
	Name           string     `gorm:"type:varchar(255);not null"`
	Slug           string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_org_app_role_slug"`
	Description    *string    `gorm:"type:text"`
	IsSystem       bool       `gorm:"not null;default:false"`
	CreatedAt      time.Time  `gorm:"not null;default:now()"`
	UpdatedAt      time.Time  `gorm:"not null;default:now()"`
}

func (Role) TableName() string {
	return "roles"
}

func RoleFromEntity(e *entity.Role) *Role {
	return &Role{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		ApplicationID:  e.ApplicationID,
		Name:           e.Name,
		Slug:           e.Slug,
		Description:    e.Description,
		IsSystem:       e.IsSystem,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
	}
}

func (m *Role) ToEntity() *entity.Role {
	return &entity.Role{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		ApplicationID:  m.ApplicationID,
		Name:           m.Name,
		Slug:           m.Slug,
		Description:    m.Description,
		IsSystem:       m.IsSystem,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
