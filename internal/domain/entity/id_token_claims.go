package entity

import "github.com/golang-jwt/jwt/v5"

type IDTokenClaims struct {
	jwt.RegisteredClaims
	Name          string `json:"name,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	AuthTime      int64  `json:"auth_time,omitempty"`
}
