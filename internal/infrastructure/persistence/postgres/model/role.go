package model

import (
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Role struct {
	ID            string     `gorm:"column:id;primaryKey"`
	ApplicationID string     `gorm:"column:application_id"`
	Code          string     `gorm:"column:code"`
	Name          string     `gorm:"column:name"`
	Description   *string    `gorm:"column:description"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
	DeletedAt     *time.Time `gorm:"column:deleted_at"`
}

func (Role) TableName() string {
	return "roles"
}

func RoleFromEntity(e *entity.Role) *Role {
	m := &Role{
		ID:        e.ID,
		Code:      e.Code,
		Name:      e.Name,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
	if e.Description != nil {
		m.Description = e.Description
	}
	if e.DeletedAt != nil && !e.DeletedAt.IsZero() {
		m.DeletedAt = e.DeletedAt
	}
	return m
}

func (m *Role) ToEntity() *entity.Role {
	e := &entity.Role{
		ID:        m.ID,
		Code:      m.Code,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Description != nil {
		e.Description = m.Description
	}
	if m.DeletedAt != nil {
		e.DeletedAt = m.DeletedAt
	}
	return e
}
