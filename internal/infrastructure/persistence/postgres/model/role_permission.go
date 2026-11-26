package model

import (
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type RolePermission struct {
	RoleID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func RolePermissionFromEntity(e *entity.RolePermission) *RolePermission {
	return &RolePermission{
		RoleID:       e.RoleID,
		PermissionID: e.PermissionID,
	}
}

func (m *RolePermission) ToEntity() *entity.RolePermission {
	return &entity.RolePermission{
		RoleID:       m.RoleID,
		PermissionID: m.PermissionID,
	}
}
