package model

import (
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Permission struct {
	ID          string     `gorm:"column:id;primaryKey"`
	Code        string     `gorm:"column:code"`
	Description *string    `gorm:"column:description"`
	Version     int        `gorm:"column:version"`
	CreatedAt   time.Time  `gorm:"column:created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

func PermissionFromEntity(e *entity.Permission) *Permission {
	m := &Permission{
		ID:          e.ID,
		Code:        e.Code,
		Description: e.Description,
		Version:     e.Version,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
	if e.Description != nil {
		m.Description = e.Description
	}
	if e.DeletedAt != nil && !e.DeletedAt.IsZero() {
		m.DeletedAt = e.DeletedAt
	}
	return m
}

func (m *Permission) ToEntity() *entity.Permission {
	e := &entity.Permission{
		ID:          m.ID,
		Code:        m.Code,
		Description: m.Description,
		Version:     m.Version,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.Description != nil {
		e.Description = m.Description
	}
	if m.DeletedAt != nil {
		e.DeletedAt = m.DeletedAt
	}
	return e
}
