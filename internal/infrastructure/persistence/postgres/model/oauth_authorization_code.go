package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type OAuthAuthorizationCode struct {
	ID                  uuid.UUID       `gorm:"type:uuid;primaryKey"`
	CodeHash            string          `gorm:"type:varchar(255);not null"`
	UserID              uuid.UUID       `gorm:"type:uuid;not null"`
	OrganizationID      *uuid.UUID      `gorm:"type:uuid"`
	ApplicationID       uuid.UUID       `gorm:"type:uuid;not null"`
	RedirectURI         string          `gorm:"type:varchar(2048);not null"`
	CodeChallenge       string          `gorm:"type:varchar(255);not null"`
	CodeChallengeMethod string          `gorm:"type:varchar(10);not null;default:'S256'"`
	Scopes              json.RawMessage `gorm:"type:jsonb;not null;default:'[]'"`
	ExpiresAt           time.Time       `gorm:"type:timestamptz;not null"`
	UsedAt              *time.Time      `gorm:"type:timestamptz"`
	CreatedAt           time.Time       `gorm:"not null;default:now()"`
}

func (OAuthAuthorizationCode) TableName() string {
	return "oauth_authorization_codes"
}

func OAuthAuthorizationCodeFromEntity(e *entity.OAuthAuthorizationCode) *OAuthAuthorizationCode {
	scopes, _ := json.Marshal(e.Scopes)

	return &OAuthAuthorizationCode{
		ID:                  e.ID,
		CodeHash:            e.CodeHash,
		UserID:              e.UserID,
		OrganizationID:      e.OrganizationID,
		ApplicationID:       e.ApplicationID,
		RedirectURI:         e.RedirectURI,
		CodeChallenge:       e.CodeChallenge,
		CodeChallengeMethod: e.CodeChallengeMethod,
		Scopes:              scopes,
		ExpiresAt:           e.ExpiresAt,
		UsedAt:              e.UsedAt,
		CreatedAt:           e.CreatedAt,
	}
}

func (m *OAuthAuthorizationCode) ToEntity() *entity.OAuthAuthorizationCode {
	var scopes []string
	_ = json.Unmarshal(m.Scopes, &scopes)

	return &entity.OAuthAuthorizationCode{
		ID:                  m.ID,
		CodeHash:            m.CodeHash,
		UserID:              m.UserID,
		OrganizationID:      m.OrganizationID,
		ApplicationID:       m.ApplicationID,
		RedirectURI:         m.RedirectURI,
		CodeChallenge:       m.CodeChallenge,
		CodeChallengeMethod: m.CodeChallengeMethod,
		Scopes:              scopes,
		ExpiresAt:           m.ExpiresAt,
		UsedAt:              m.UsedAt,
		CreatedAt:           m.CreatedAt,
	}
}
