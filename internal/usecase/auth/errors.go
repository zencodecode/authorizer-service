package auth

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidClient            = errors.New("invalid client_id")
	ErrRedirectURINotRegistered = errors.New("redirect_uri not registered for this client")
	ErrInvalidCredentials       = errors.New("email or password is invalid")
	ErrEmailNotVerified         = errors.New("email is not verified")
	ErrAccountSuspended         = errors.New("account is suspended")
)

type AuthError struct {
	Code        string
	Description string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Description)
}

func newAuthError(code, description string) *AuthError {
	return &AuthError{Code: code, Description: description}
}
