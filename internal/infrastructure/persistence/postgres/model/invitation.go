package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Invitation struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	OrganizationID uuid.UUID  `gorm:"type:uuid;not null"`
	Email          string     `gorm:"type:varchar(255);not null"`
	RoleID         *uuid.UUID `gorm:"type:uuid"`
	InvitedBy      uuid.UUID  `gorm:"type:uuid;not null"`
	Token          string     `gorm:"type:varchar(255);not null;uniqueIndex"`
	Status         string     `gorm:"type:varchar(50);not null;default:'pending'"`
	ExpiresAt      time.Time  `gorm:"type:timestamptz;not null"`
	CreatedAt      time.Time  `gorm:"not null;default:now()"`
}

func (Invitation) TableName() string {
	return "invitations"
}

func InvitationFromEntity(e *entity.Invitation) *Invitation {
	return &Invitation{
		ID:             e.ID,
		OrganizationID: e.OrganizationID,
		Email:          e.Email,
		RoleID:         e.RoleID,
		InvitedBy:      e.InvitedBy,
		Token:          e.Token,
		Status:         e.Status,
		ExpiresAt:      e.ExpiresAt,
		CreatedAt:      e.CreatedAt,
	}
}

func (m *Invitation) ToEntity() *entity.Invitation {
	return &entity.Invitation{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		Email:          m.Email,
		RoleID:         m.RoleID,
		InvitedBy:      m.InvitedBy,
		Token:          m.Token,
		Status:         m.Status,
		ExpiresAt:      m.ExpiresAt,
		CreatedAt:      m.CreatedAt,
	}
}
