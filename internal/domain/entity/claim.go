package entity

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Name        string
	Email       string
	OrgID       *string
	OrgSlug     *string
	Scopes      []string
	Roles       []string
	Permissions []string
}
