package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Application struct {
	ID          uuid.UUID      `gorm:"column:id;primaryKey"`
	Code        string         `gorm:"column:code;"`
	Name        string         `gorm:"column:name;"`
	Description *string        `gorm:"column:description;"`
	Metadata    map[string]any `gorm:"column:metadata;"`
	CreatedAt   time.Time      `gorm:"column:created_at;"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;"`
	DeletedAt   *time.Time     `gorm:"column:deleted_at;"`
}

func (Application) TableName() string {
	return "applications"
}

func ApplicationFromEntity(e *entity.Application) *Application {
	m := &Application{
		ID:          e.ID,
		Code:        e.Code,
		Name:        e.Name,
		Description: e.Description,
		Metadata:    e.Metadata,
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

func (m *Application) ToEntity() *entity.Application {
	e := &entity.Application{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Metadata:    m.Metadata,
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
