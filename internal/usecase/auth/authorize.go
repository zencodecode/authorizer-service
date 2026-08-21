package auth

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/applicationscope"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/authorizesession"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/randutil"
)

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
		LoginChallengeID     string
	}
)

type authorizeUsecase struct {
	appRepo     application.Repository
	scopeRepo   applicationscope.Repository
	sessionRepo authorizesession.Repository
	logger      service.Logger
}

func NewAuthorizeUsecase(
	appRepo application.Repository,
	scopeRepo applicationscope.Repository,
	sessionRepo authorizesession.Repository,
	logger service.Logger,
) AuthorizeUsecase {
	return &authorizeUsecase{
		appRepo:     appRepo,
		scopeRepo:   scopeRepo,
		sessionRepo: sessionRepo,
		logger:      logger,
	}
}

func (uc *authorizeUsecase) Execute(ctx context.Context, params AuthorizeParams) (*AuthorizeResult, error) {
	if params.ResponseType != "code" {
		return nil, newAuthError("unsupported_response_type", "only 'code' response_type is supported")
	}

	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query application",
			"action", "AUTHORIZE",
			"client_id", params.ClientID,
			"error", err.Error(),
		)
		return nil, ErrInvalidClient
	}

	if !isRedirectURIAllowed(app.RedirectURIs, params.RedirectURI) {
		uc.logger.Warn(ctx, "redirect_uri not registered",
			"action", "AUTHORIZE",
			"client_id", params.ClientID,
			"redirect_uri", params.RedirectURI,
		)
		return nil, ErrRedirectURINotRegistered
	}

	if params.CodeChallengeMethod != "S256" {
		return nil, newAuthError("invalid_request", "code_challenge_method must be S256")
	}
	if len(params.CodeChallenge) < 43 {
		return nil, newAuthError("invalid_request", "code_challenge is required and must be at least 43 characters")
	}

	if len(params.Scope) == 0 {
		return nil, newAuthError("invalid_scope", "at least one scope is required")
	}
	for _, scope := range params.Scope {
		s, err := uc.scopeRepo.GetByApplicationAndScope(ctx, app.ID, scope)
		if err != nil {
			uc.logger.Error(ctx, "failed to query scope",
				"action", "AUTHORIZE",
				"scope", scope,
				"error", err.Error(),
			)
			return nil, newAuthError("server_error", "failed to validate scopes")
		}
		if s == nil {
			return nil, newAuthError("invalid_scope", fmt.Sprintf("scope %q is not registered for this application", scope))
		}
	}

	challengeID, err := randutil.GenerateRandomString(32)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate challenge id",
			"action", "AUTHORIZE",
			"client_id", params.ClientID,
			"error", err.Error())
		return nil, newAuthError("server_error", "failed to generate challenge id")
	}

	sess := entity.AuthorizeSession{
		ClientID:            params.ClientID,
		RedirectURI:         params.RedirectURI,
		Scope:               params.Scope,
		State:               params.State,
		CodeChallenge:       params.CodeChallenge,
		CodeChallengeMethod: params.CodeChallengeMethod,
	}

	if err := uc.sessionRepo.Save(ctx, challengeID, sess, 10*time.Minute); err != nil {
		uc.logger.Error(ctx, "failed to save authorize session",
			"action", "AUTHORIZE",
			"client_id", params.ClientID,
			"error", err.Error())
		return nil, newAuthError("server_error", "failed to persist authorize request")
	}

	return &AuthorizeResult{
		ClientID:             app.ClientID,
		RedirectURI:          params.RedirectURI,
		Scopes:               params.Scope,
		State:                params.State,
		CodeChallenge:        params.CodeChallenge,
		CodeChallengeMethod:  params.CodeChallengeMethod,
		RequiresOrganization: app.RequiresOrganization,
		ApplicationName:      app.Name,
		LoginChallengeID:     challengeID,
	}, nil
}

func isRedirectURIAllowed(registered []string, requested string) bool {
	return slices.Contains(registered, requested)
}

// ScopeString returns scopes as a space-separated string.
func (r *AuthorizeResult) ScopeString() string {
	return strings.Join(r.Scopes, " ")
}
