package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type User struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Email           string     `gorm:"type:citext;not null;uniqueIndex"`
	PasswordHash    string     `gorm:"type:varchar(255);not null"`
	Name            string     `gorm:"type:varchar(255);not null"`
	Status          string     `gorm:"type:varchar(50);not null;default:'pending_verification'"`
	EmailVerifiedAt *time.Time `gorm:"type:timestamptz"`
	CreatedAt       time.Time  `gorm:"not null;default:now()"`
	UpdatedAt       time.Time  `gorm:"not null;default:now()"`
	DeletedAt       *time.Time `gorm:"index"`
}

func (User) TableName() string {
	return "users"
}

func UserFromEntity(e *entity.User) *User {
	return &User{
		ID:              e.ID,
		Email:           e.Email,
		PasswordHash:    e.PasswordHash,
		Name:            e.Name,
		Status:          e.Status,
		EmailVerifiedAt: e.EmailVerifiedAt,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		DeletedAt:       e.DeletedAt,
	}
}

func (m *User) ToEntity() *entity.User {
	return &entity.User{
		ID:              m.ID,
		Email:           m.Email,
		PasswordHash:    m.PasswordHash,
		Name:            m.Name,
		Status:          m.Status,
		EmailVerifiedAt: m.EmailVerifiedAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt:       m.DeletedAt,
	}
}
