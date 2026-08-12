package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Application struct {
	ID                   uuid.UUID       `gorm:"type:uuid;primaryKey"`
	Name                 string          `gorm:"type:varchar(255);not null"`
	Slug                 string          `gorm:"type:varchar(255);not null;uniqueIndex"`
	ClientID             string          `gorm:"type:varchar(255);not null;uniqueIndex"`
	ClientSecretHash     string          `gorm:"type:varchar(255);not null"`
	RedirectURIs         json.RawMessage `gorm:"type:jsonb;not null;default:'[]'"`
	AllowedGrantTypes    json.RawMessage `gorm:"type:jsonb;not null;default:'[\"authorization_code\",\"refresh_token\"]'"`
	RequiresOrganization bool            `gorm:"not null;default:true"`
	Metadata             json.RawMessage `gorm:"type:jsonb"`
	IsActive             bool            `gorm:"not null;default:true"`
	CreatedAt            time.Time       `gorm:"not null;default:now()"`
	UpdatedAt            time.Time       `gorm:"not null;default:now()"`
}

func (Application) TableName() string {
	return "applications"
}

func ApplicationFromEntity(e *entity.Application) *Application {
	redirectURIs, _ := json.Marshal(e.RedirectURIs)
	allowedGrantTypes, _ := json.Marshal(e.AllowedGrantTypes)

	var metadata json.RawMessage
	if e.Metadata != nil {
		metadata, _ = json.Marshal(e.Metadata)
	}

	return &Application{
		ID:                   e.ID,
		Name:                 e.Name,
		Slug:                 e.Slug,
		ClientID:             e.ClientID,
		ClientSecretHash:     e.ClientSecretHash,
		RedirectURIs:         redirectURIs,
		AllowedGrantTypes:    allowedGrantTypes,
		RequiresOrganization: e.RequiresOrganization,
		Metadata:             metadata,
		IsActive:             e.IsActive,
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
	}
}

func (m *Application) ToEntity() *entity.Application {
	var redirectURIs []string
	_ = json.Unmarshal(m.RedirectURIs, &redirectURIs)

	var allowedGrantTypes []string
	_ = json.Unmarshal(m.AllowedGrantTypes, &allowedGrantTypes)

	var metadata map[string]any
	if m.Metadata != nil {
		_ = json.Unmarshal(m.Metadata, &metadata)
	}

	return &entity.Application{
		ID:                   m.ID,
		Name:                 m.Name,
		Slug:                 m.Slug,
		ClientID:             m.ClientID,
		ClientSecretHash:     m.ClientSecretHash,
		RedirectURIs:         redirectURIs,
		AllowedGrantTypes:    allowedGrantTypes,
		RequiresOrganization: m.RequiresOrganization,
		Metadata:             metadata,
		IsActive:             m.IsActive,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}
