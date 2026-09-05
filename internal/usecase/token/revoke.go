package token

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
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
	oauthRefRepo oauthrefreshtoken.Repository
	logger       service.Logger
}

func NewRevokeUsecase(
	appRepo application.Repository,
	oauthRefRepo oauthrefreshtoken.Repository,
	logger service.Logger,
) RevokeUsecase {
	return &revokeUsecase{
		appRepo:      appRepo,
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

	if !hash.CheckHash(app.ClientSecretHash, params.ClientSecret) {
		return apperr.NewDirectError(enum.INVALID_CLIENT, "invalid client credentials")
	}

	tokenHash := hash.HashSHA256(params.Token)

	refToken, err := uc.oauthRefRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		uc.logger.Warn(ctx, "revoke requested for unknown token",
			"action", "REVOKE",
			"client_id", params.ClientID,
		)
		return nil
	}

	if err := uc.oauthRefRepo.Revoke(ctx, refToken.ID); err != nil {
		uc.logger.Error(ctx, "failed to revoke refresh token",
			"action", "REVOKE",
			"error", err.Error())
		return apperr.NewDirectError(enum.SERVER_ERROR, "failed to revoke token")
	}

	return nil
}
