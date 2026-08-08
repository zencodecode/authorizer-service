package entity

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	jwt.RegisteredClaims
	Username      string          `json:"username"`
	Email         string          `json:"email"`
	Authorization []Authorization `json:"authorization"`
}

type Authorization struct {
	Application AuthorizationApplication `json:"application"`
}

type AuthorizationApplication struct {
	Code  string              `json:"code"`
	Name  string              `json:"name"`
	Roles []AuthorizationRole `json:"roles"`
}

type AuthorizationRole struct {
	Code        string                    `json:"code"`
	Name        string                    `json:"name"`
	Permissions []AuthorizationPermission `json:"permissions"`
}

type AuthorizationPermission struct {
	Code string `json:"code"`
	Name string `json:"name"`
}
