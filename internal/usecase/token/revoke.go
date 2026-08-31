package token

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthaccesstoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthrefreshtoken"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	RevokeParams struct {
		Token         string
		TokenTypeHint string
		ClientID      string
		ClientSecret  string
	}
)

type revokeUsecase struct {
	appRepo      application.Repository
	oauthAccRepo oauthaccesstoken.Repository
	oauthRefRepo oauthrefreshtoken.Repository
	logger       service.Logger
}

func NewRevokeUsecase(
	appRepo application.Repository,
	oauthAccRepo oauthaccesstoken.Repository,
	oauthRefRepo oauthrefreshtoken.Repository,
	logger service.Logger,
) RevokeUsecase {
	return &revokeUsecase{
		appRepo:      appRepo,
		oauthAccRepo: oauthAccRepo,
		oauthRefRepo: oauthRefRepo,
		logger:       logger,
	}
}

func (uc *revokeUsecase) Execute(ctx context.Context, params RevokeParams) error {
	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query applicaton",
			"action", "REVOKE",
			"client_id", &params.ClientID,
			"error", err.Error(),
		)
		return apperr.NewDirectError(enum.SERVER_ERROR, "failed to query application")
	}

	tokenHash := hash.HashSHA256(params.Token)

	primary, fallback := uc.revokeRefreshToken, uc.revokeAccessToken
	if params.TokenTypeHint != "refresh_token" {
		primary, fallback = uc.revokeAccessToken, uc.revokeRefreshToken
	}

	revoked, err := primary(ctx, app.ID, tokenHash)
	if err != nil {
		return err
	}
	if revoked {
		return nil
	}

	revoked, err = fallback(ctx, app.ID, tokenHash)
	if err != nil {
		return err
	}
	if !revoked {
		uc.logger.Warn(ctx, "revoke requested for unknown token",
			"action", "REVOKE", "client_id", params.ClientID)
	}

	return nil
}

func (uc *revokeUsecase) revokeRefreshToken(ctx context.Context, appID uuid.UUID, tokenHash string) (bool, error) {
	refToken, err := uc.oauthRefRepo.GetByTokenHash(ctx, tokenHash)
	if errors.Is(err, oauthrefreshtoken.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		uc.logger.Error(ctx, "failed to query oauth refresh token",
			"action", "REVOKE", "application_id", appID, "error", err.Error())
		return false, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query oauth refresh token")
	}

	accToken, err := uc.oauthAccRepo.GetByID(ctx, refToken.AccessTokenID)
	if errors.Is(err, oauthrefreshtoken.ErrNotFound) {
		uc.logger.Warn(ctx, "refresh token references missing access token",
			"action", "REVOKE", "access_token_id", refToken.AccessTokenID)
		return false, nil
	}
	if err != nil {
		uc.logger.Error(ctx, "failed to query oauth access token",
			"action", "REVOKE", "application_id", appID, "error", err.Error())
		return false, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query oauth access token")
	}
	if accToken.ApplicationID != appID {
		uc.logger.Warn(ctx, "client attempted to revoke refresh token belonging to another application",
			"action", "REVOKE",
			"application_id", appID,
			"token_application_id", accToken.ApplicationID,
		)
		return false, nil
	}

	if err := uc.oauthRefRepo.Revoke(ctx, refToken.ID); err != nil {
		uc.logger.Error(ctx, "failed to revoke refresh token",
			"action", "REVOKE", "application_id", appID, "error", err.Error())
		return false, apperr.NewDirectError(enum.SERVER_ERROR, "failed to revoke refresh token")
	}
	return true, nil
}

func (uc *revokeUsecase) revokeAccessToken(ctx context.Context, appID uuid.UUID, tokenHash string) (bool, error) {
	token, err := uc.oauthAccRepo.GetByTokenHash(ctx, tokenHash)
	if errors.Is(err, oauthaccesstoken.ErrNotFound) {
		return false, nil
	}
	if err != nil {
		uc.logger.Error(ctx, "failed to query oauth access token",
			"action", "REVOKE", "application_id", appID, "error", err.Error())
		return false, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query oauth access token")
	}

	if token.ApplicationID != appID {
		uc.logger.Warn(ctx, "client attempted to revoke refresh token belonging to another application",
			"action", "REVOKE",
			"application_id", appID,
			"token_application_id", token.ApplicationID,
		)
		return false, nil
	}

	if err := uc.oauthAccRepo.Revoke(ctx, token.ID); err != nil {
		uc.logger.Error(ctx, "failed to revoke access token",
			"action", "REVOKE", "application_id", appID, "error", err.Error())
		return false, apperr.NewDirectError(enum.SERVER_ERROR, "failed to revoke access token")
	}
	return true, nil
}
