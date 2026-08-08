package entity

import "time"

type User struct {
	ID            string
	Email         string
	FullName      string
	Password      string
	Phone         *string
	IsActive      bool
	EmailVerified bool
	PhoneVerified bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
