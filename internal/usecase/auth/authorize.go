package auth

import (
	"context"
	"errors"
	"slices"

	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/applicationscope"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
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

	AuthorizeOutput struct {
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

func (uc *authorizeUsecase) Execute(ctx context.Context, params AuthorizeParams) (*AuthorizeOutput, error) {
	if params.ResponseType != "code" {
		return nil, errors.New("unsupported response type")
	}

	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Warn(ctx, "authorize failed: invalid client_id",
			"client_id", params.ClientID,
		)
		return nil, errors.New("client not found")
	}

	if isRedirectURIAllowed(app.RedirectURIs, params.RedirectURI) {
		uc.logger.Warn(ctx, "authorize failed: redirect_uri not registered for this client",
			"client_id", params.ClientID,
		)
		return nil, errors.New("redirect_uri not registered for this client")
	}

	for _, scope := range params.Scope {
		_, err := uc.scopeRepo.GetByApplicationAndScope(ctx, app.ID, scope)
		if err != nil {
			uc.logger.Warn(ctx, "authorize failed: invalid scope",
				"client_id", params.ClientID,
			)
			return nil, errors.New("scope not found")
		}
	}

	if params.CodeChallengeMethod != "S256" {
		return nil, errors.New("invalid challenge method")
	}

	if params.CodeChallenge == "" || len(params.CodeChallenge) < 43 {
		return nil, errors.New("invalid code challenge")
	}

	return nil, nil
}

func isRedirectURIAllowed(registered []string, requested string) bool {
	return slices.Contains(registered, requested)
}
