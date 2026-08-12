package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type UserRole struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	OrganizationID *uuid.UUID `gorm:"type:uuid;uniqueIndex:idx_org_user_role"`
	UserID         uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_org_user_role"`
	RoleID         uuid.UUID  `gorm:"type:uuid;not null;uniqueIndex:idx_org_user_role"`
	AssignedAt     time.Time  `gorm:"type:timestamptz;not null;default:now()"`
	AssignedBy     *uuid.UUID `gorm:"type:uuid"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

func UserRoleFromEntity(e *entity.UserRole) *UserRole {
	return &UserRole{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		UserID:         e.UserID,
		RoleID:         e.RoleID,
		AssignedAt:     e.AssignedAt,
		AssignedBy:     e.AssignedBy,
	}
}

func (m *UserRole) ToEntity() *entity.UserRole {
	return &entity.UserRole{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		RoleID:         m.RoleID,
		AssignedAt:     m.AssignedAt,
		AssignedBy:     m.AssignedBy,
	}
}
