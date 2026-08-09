package model

import (
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type UserRole struct {
	UserID    string    `gorm:"column:user_id"`
	RoleID    string    `gorm:"column:role_id"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

func UserRoleFromEntity(e *entity.UserRole) *UserRole {
	m := &UserRole{
		UserID:    e.UserID,
		RoleID:    e.RoleID,
		CreatedAt: e.CreatedAt,
	}
	return m
}

func (m *UserRole) ToEntity() *entity.UserRole {
	e := &entity.UserRole{
		UserID:    m.UserID,
		RoleID:    m.RoleID,
		CreatedAt: m.CreatedAt,
	}
	return e
}
