package apperr

import (
	"fmt"

	"github.com/zencodecode/authorizer-service/internal/definition/enum"
)

type OAuthError struct {
	Code        enum.OAuthError
	Description string
	RedirectURI string
	State       string
}

func (e *OAuthError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Description)
}

func NewRedirectableError(code enum.OAuthError, description, redirectURI, state string) *OAuthError {
	return &OAuthError{Code: code, Description: description, RedirectURI: redirectURI, State: state}
}

func NewFatalError(code enum.OAuthError, description string) *OAuthError {
	return &OAuthError{Code: code, Description: description}
}
