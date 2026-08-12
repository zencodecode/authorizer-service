package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Permission struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	ApplicationID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_app_perm_slug"`
	Slug          string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_app_perm_slug"`
	Resource      string    `gorm:"type:varchar(255);not null"`
	Action        string    `gorm:"type:varchar(255);not null"`
	Description   *string   `gorm:"type:text"`
	CreatedAt     time.Time `gorm:"not null;default:now()"`
}

func (Permission) TableName() string {
	return "permissions"
}

func PermissionFromEntity(e *entity.Permission) *Permission {
	return &Permission{
		ID:            e.ID,
		ApplicationID: e.ApplicationID,
		Slug:          e.Slug,
		Resource:      e.Resource,
		Action:        e.Action,
		Description:   e.Description,
		CreatedAt:     e.CreatedAt,
	}
}

func (m *Permission) ToEntity() *entity.Permission {
	return &entity.Permission{
		ID:            m.ID,
		ApplicationID: m.ApplicationID,
		Slug:          m.Slug,
		Resource:      m.Resource,
		Action:        m.Action,
		Description:   m.Description,
		CreatedAt:     m.CreatedAt,
	}
}
