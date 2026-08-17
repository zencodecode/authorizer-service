package entity

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	jwt.RegisteredClaims
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	OrgID   *string  `json:"org_id,omitempty"`
	OrgSlug *string  `json:"org_slug,omitempty"`
	Scopes  []string `json:"scopes"`

	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}
