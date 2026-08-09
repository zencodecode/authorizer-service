package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type User struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email         string    `gorm:"type:citext;uniqueIndex;not null"`
	Username      string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password      string    `gorm:"type:text;not null"`
	FullName      string    `gorm:"type:varchar(100)"`
	Phone         *string   `gorm:"type:varchar(20)"`
	IsActive      bool      `gorm:"type:bool;default:true;not null"`
	EmailVerified bool      `gorm:"type:bool;default:false;not null"`
	PhoneVerified bool      `gorm:"type:bool;default:false;not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}

func (User) TableName() string {
	return "users"
}

func UserFromEntity(e *entity.User) *User {
	m := &User{
		ID:            e.ID,
		Email:         e.Email,
		Username:      e.Username,
		Password:      e.Password,
		FullName:      e.FullName,
		Phone:         e.Phone,
		IsActive:      e.IsActive,
		EmailVerified: e.EmailVerified,
		PhoneVerified: e.PhoneVerified,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
	}
	if e.DeletedAt != nil && !e.DeletedAt.IsZero() {
		m.DeletedAt = e.DeletedAt
	}
	return m
}

func (m *User) ToEntity() *entity.User {
	e := &entity.User{
		ID:            m.ID,
		Email:         m.Email,
		Username:      m.Username,
		Password:      m.Password,
		FullName:      m.FullName,
		Phone:         m.Phone,
		IsActive:      m.IsActive,
		EmailVerified: m.EmailVerified,
		PhoneVerified: m.PhoneVerified,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
	}
	if m.DeletedAt != nil {
		e.DeletedAt = m.DeletedAt
	}
	return e
}
