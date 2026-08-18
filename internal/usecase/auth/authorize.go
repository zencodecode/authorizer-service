package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/applicationscope"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

var (
	ErrInvalidClient            = errors.New("invalid client_id")
	ErrRedirectURINotRegistered = errors.New("redirect_uri not registered for this client")
)

type AuthorizeError struct {
	Code        string
	Description string
}

func (e *AuthorizeError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Description)
}

func newAuthorizeError(code, description string) *AuthorizeError {
	return &AuthorizeError{Code: code, Description: description}
}

type (
	AuthorizeParams struct {
		ResponseType        string
		ClientID            string
		RedirectURI         string
		Scope               []string
		State               string
		CodeChallenge       string
		CodeChallengeMethod string
	}

	AuthorizeResult struct {
		ClientID             string
		RedirectURI          string
		Scopes               []string
		State                string
		CodeChallenge        string
		CodeChallengeMethod  string
		RequiresOrganization bool
		ApplicationName      string
	}
)

type authorizeUsecase struct {
	appRepo   application.Repository
	scopeRepo applicationscope.Repository
	logger    service.Logger
}

func NewAuthorizeUsecase(
	appRepo application.Repository,
	scopeRepo applicationscope.Repository,
	logger service.Logger,
) AuthorizeUsecase {
	return &authorizeUsecase{
		appRepo:   appRepo,
		scopeRepo: scopeRepo,
		logger:    logger,
	}
}

func (uc *authorizeUsecase) Execute(ctx context.Context, params AuthorizeParams) (*AuthorizeResult, error) {
	// 1. Validate response_type (check early before any DB call)
	if params.ResponseType != "code" {
		return nil, newAuthorizeError("unsupported_response_type", "only 'code' response_type is supported")
	}

	// 2. Lookup and validate application
	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "authorize: failed to query application",
			"client_id", params.ClientID,
			"error", err.Error(),
		)
		return nil, ErrInvalidClient
	}
	if app == nil || !app.IsActive {
		uc.logger.Warn(ctx, "authorize: invalid or inactive client",
			"client_id", params.ClientID,
		)
		return nil, ErrInvalidClient
	}

	// 3. Validate redirect_uri (MUST check before redirecting any error)
	if !isRedirectURIAllowed(app.RedirectURIs, params.RedirectURI) {
		uc.logger.Warn(ctx, "authorize: redirect_uri not registered",
			"client_id", params.ClientID,
			"redirect_uri", params.RedirectURI,
		)
		return nil, ErrRedirectURINotRegistered
	}

	// From here, errors can be safely redirected back to client

	// 4. Validate PKCE
	if params.CodeChallengeMethod != "S256" {
		return nil, newAuthorizeError("invalid_request", "code_challenge_method must be S256")
	}
	if len(params.CodeChallenge) < 43 {
		return nil, newAuthorizeError("invalid_request", "code_challenge is required and must be at least 43 characters")
	}

	// 5. Validate scopes
	if len(params.Scope) == 0 {
		return nil, newAuthorizeError("invalid_scope", "at least one scope is required")
	}
	for _, scope := range params.Scope {
		s, err := uc.scopeRepo.GetByApplicationAndScope(ctx, app.ID, scope)
		if err != nil {
			uc.logger.Error(ctx, "authorize: failed to query scope",
				"scope", scope,
				"error", err.Error(),
			)
			return nil, newAuthorizeError("server_error", "failed to validate scopes")
		}
		if s == nil {
			return nil, newAuthorizeError("invalid_scope", fmt.Sprintf("scope %q is not registered for this application", scope))
		}
	}

	// 6. Validation passed — return result for handler to proceed (show login page)
	return &AuthorizeResult{
		ClientID:             app.ClientID,
		RedirectURI:          params.RedirectURI,
		Scopes:               params.Scope,
		State:                params.State,
		CodeChallenge:        params.CodeChallenge,
		CodeChallengeMethod:  params.CodeChallengeMethod,
		RequiresOrganization: app.RequiresOrganization,
		ApplicationName:      app.Name,
	}, nil
}

func isRedirectURIAllowed(registered []string, requested string) bool {
	return slices.Contains(registered, requested)
}

// ScopeString returns scopes as a space-separated string.
func (r *AuthorizeResult) ScopeString() string {
	return strings.Join(r.Scopes, " ")
}
