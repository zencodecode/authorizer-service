package entity

import "github.com/google/uuid"

type AuthorizeSession struct {
	ClientID            string
	RedirectURI         string
	Scope               []string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	UserID              *uuid.UUID
}
