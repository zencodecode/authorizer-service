package oauthauthorizationcode

import "errors"

var (
	ErrNotFound          = errors.New("authorization code not found")
	ErrExpired           = errors.New("authorization code expired")
	ErrExpiredOrNotFound = errors.New("authorization code expired or not found")
)
