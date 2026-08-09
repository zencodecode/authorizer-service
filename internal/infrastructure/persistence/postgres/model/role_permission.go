package model

import (
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type RolePermission struct {
	RoleID       string    `gorm:"column:role_id"`
	PermissionID string    `gorm:"column:permission_id"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

func RolePermissionFromEntity(e *entity.RolePermission) *RolePermission {
	m := &RolePermission{
		RoleID:       e.RoleID,
		PermissionID: e.PermissionID,
		CreatedAt:    e.CreatedAt,
	}
	return m
}

func (m *RolePermission) ToEntity() *entity.RolePermission {
	e := &entity.RolePermission{
		RoleID:       m.RoleID,
		PermissionID: m.PermissionID,
		CreatedAt:    m.CreatedAt,
	}
	return e
}
