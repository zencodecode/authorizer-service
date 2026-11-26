package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Organization struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Name      string     `gorm:"type:varchar(255);not null"`
	Slug      string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	Status    string     `gorm:"type:varchar(50);not null;default:'active'"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
	UpdatedAt time.Time  `gorm:"not null;default:now()"`
	DeletedAt *time.Time `gorm:"index"`
}

func (Organization) TableName() string {
	return "organizations"
}

func OrganizationFromEntity(e *entity.Organization) *Organization {
	return &Organization{
		ID:        e.ID,
		Name:      e.Name,
		Slug:      e.Slug,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,
	}
}

func (m *Organization) ToEntity() *entity.Organization {
	return &entity.Organization{
		ID:        m.ID,
		Name:      m.Name,
		Slug:      m.Slug,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}
