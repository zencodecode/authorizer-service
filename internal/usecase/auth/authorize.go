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
		ResponseType        string
		ClientID            string
		RedirectURI         string
		Scope               string
		State               string
		CodeChallenge       string
		CodeChallengeMethod string
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
	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Warn(ctx, "authorize failed: invalid client_id",
			"client_id", params.ClientID,
		)
		return nil, ErrInvalidClient
	}

	if !isRedirectURIAllowed(app.RedirectURIs, params.RedirectURI) {
		uc.logger.Warn(ctx, "authorize failed: redirect_uri not registered for this client",
			"client_id", params.ClientID,
		)
		return nil, ErrRedirectURINotRegistered
	}

	if params.ResponseType != "code" {
		return nil, newAuthorizeError("unsupported_response_type", "only 'code' response_type is supported")
	}

	for _, scope := range params.Scope {
		_, err := uc.scopeRepo.GetByApplicationAndScope(ctx, app.ID, scope)
		if err != nil {
			uc.logger.Warn(ctx, "authorize failed: invalid scope",
				"client_id", params.ClientID,
				"scope", scope,
			)
			return nil, newAuthorizeError("invalid_scope", fmt.Sprintf("scope %q is not registered for this client", scope))
		}
	}

	if params.CodeChallengeMethod != "S256" {
		return nil, newAuthorizeError("invalid_request", "code_challenge_method must be S256")
	}

	if params.CodeChallenge == "" || len(params.CodeChallenge) < 43 {
		return nil, newAuthorizeError("invalid_request", "code_challenge is missing or too short")
	}

	return &AuthorizeResult{
		ResponseType:        params.ResponseType,
		ClientID:            params.ClientID,
		RedirectURI:         params.RedirectURI,
		Scope:               strings.Join(params.Scope, " "),
		State:               params.State,
		CodeChallenge:       params.CodeChallenge,
		CodeChallengeMethod: params.CodeChallengeMethod,
	}, nil
}

func isRedirectURIAllowed(registered []string, requested string) bool {
	return slices.Contains(registered, requested)
}
