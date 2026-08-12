package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type OAuthAccessToken struct {
	ID             uuid.UUID       `gorm:"type:uuid;primaryKey"`
	TokenHash      string          `gorm:"type:varchar(255);not null"`
	UserID         uuid.UUID       `gorm:"type:uuid;not null"`
	OrganizationID *uuid.UUID      `gorm:"type:uuid"`
	ApplicationID  uuid.UUID       `gorm:"type:uuid;not null"`
	Scopes         json.RawMessage `gorm:"type:jsonb;not null;default:'[]'"`
	ExpiresAt      time.Time       `gorm:"type:timestamptz;not null"`
	RevokedAt      *time.Time      `gorm:"type:timestamptz"`
	CreatedAt      time.Time       `gorm:"not null;default:now()"`
}

func (OAuthAccessToken) TableName() string {
	return "oauth_access_tokens"
}

func OAuthAccessTokenFromEntity(e *entity.OAuthAccessToken) *OAuthAccessToken {
	scopes, _ := json.Marshal(e.Scopes)

	return &OAuthAccessToken{
		ID:             e.ID,
		TokenHash:      e.TokenHash,
		UserID:         e.UserID,
		OrganizationID: e.OrganizationID,
		ApplicationID:  e.ApplicationID,
		Scopes:         scopes,
		ExpiresAt:      e.ExpiresAt,
		RevokedAt:      e.RevokedAt,
		CreatedAt:      e.CreatedAt,
	}
}

func (m *OAuthAccessToken) ToEntity() *entity.OAuthAccessToken {
	var scopes []string
	_ = json.Unmarshal(m.Scopes, &scopes)

	return &entity.OAuthAccessToken{
		ID:             m.ID,
		TokenHash:      m.TokenHash,
		UserID:         m.UserID,
		OrganizationID: m.OrganizationID,
		ApplicationID:  m.ApplicationID,
		Scopes:         scopes,
		ExpiresAt:      m.ExpiresAt,
		RevokedAt:      m.RevokedAt,
		CreatedAt:      m.CreatedAt,
	}
}
