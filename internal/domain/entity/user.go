package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID
	Email         string
	Username      string
	Password      string
	FullName      string
	Phone         *string
	IsActive      bool
	EmailVerified bool
	PhoneVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
