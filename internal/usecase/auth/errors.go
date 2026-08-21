package auth

import "errors"

var (
	ErrInvalidCredentials = errors.New("email or password is invalid")
	ErrEmailNotVerified   = errors.New("email is not verified")
	ErrAccountSuspended   = errors.New("account is suspended")
)
