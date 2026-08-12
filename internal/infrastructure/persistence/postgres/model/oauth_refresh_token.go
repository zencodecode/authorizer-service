package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type OAuthRefreshToken struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	AccessTokenID uuid.UUID  `gorm:"type:uuid;not null"`
	TokenHash     string     `gorm:"type:varchar(255);not null"`
	UserAgent     *string    `gorm:"type:varchar(512)"`
	IPAddress     *string    `gorm:"type:varchar(45)"`
	ExpiresAt     time.Time  `gorm:"type:timestamptz;not null"`
	RevokedAt     *time.Time `gorm:"type:timestamptz"`
	CreatedAt     time.Time  `gorm:"not null;default:now()"`
}

func (OAuthRefreshToken) TableName() string {
	return "oauth_refresh_tokens"
}

func OAuthRefreshTokenFromEntity(e *entity.OAuthRefreshToken) *OAuthRefreshToken {
	return &OAuthRefreshToken{
		ID:            e.ID,
		AccessTokenID: e.AccessTokenID,
		TokenHash:     e.TokenHash,
		UserAgent:     e.UserAgent,
		IPAddress:     e.IPAddress,
		ExpiresAt:     e.ExpiresAt,
		RevokedAt:     e.RevokedAt,
		CreatedAt:     e.CreatedAt,
	}
}

func (m *OAuthRefreshToken) ToEntity() *entity.OAuthRefreshToken {
	return &entity.OAuthRefreshToken{
		ID:            m.ID,
		AccessTokenID: m.AccessTokenID,
		TokenHash:     m.TokenHash,
		UserAgent:     m.UserAgent,
		IPAddress:     m.IPAddress,
		ExpiresAt:     m.ExpiresAt,
		RevokedAt:     m.RevokedAt,
		CreatedAt:     m.CreatedAt,
	}
}
