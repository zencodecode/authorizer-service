package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type PasswordResetToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	TokenHash string     `gorm:"type:varchar(255);not null"`
	ExpiresAt time.Time  `gorm:"type:timestamptz;not null"`
	UsedAt    *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

func PasswordResetTokenFromEntity(e *entity.PasswordResetToken) *PasswordResetToken {
	return &PasswordResetToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		ExpiresAt: e.ExpiresAt,
		UsedAt:    e.UsedAt,
		CreatedAt: e.CreatedAt,
	}
}

func (m *PasswordResetToken) ToEntity() *entity.PasswordResetToken {
	return &entity.PasswordResetToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
		CreatedAt: m.CreatedAt,
	}
}

type EmailVerificationToken struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	TokenHash string     `gorm:"type:varchar(255);not null"`
	ExpiresAt time.Time  `gorm:"type:timestamptz;not null"`
	UsedAt    *time.Time `gorm:"type:timestamptz"`
	CreatedAt time.Time  `gorm:"not null;default:now()"`
}

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

func EmailVerificationTokenFromEntity(e *entity.EmailVerificationToken) *EmailVerificationToken {
	return &EmailVerificationToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		ExpiresAt: e.ExpiresAt,
		UsedAt:    e.UsedAt,
		CreatedAt: e.CreatedAt,
	}
}

func (m *EmailVerificationToken) ToEntity() *entity.EmailVerificationToken {
	return &entity.EmailVerificationToken{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
		CreatedAt: m.CreatedAt,
	}
}
