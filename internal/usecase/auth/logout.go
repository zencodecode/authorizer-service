package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthaccesstoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthrefreshtoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type (
	LogoutParams struct {
		UserID   uuid.UUID
		ClientID string
	}
)

type logoutUsecase struct {
	userRepo     user.Repository
	appRepo      application.Repository
	oauthAccRepo oauthaccesstoken.Repository
	oauthRefRepo oauthrefreshtoken.Repository
	logger       service.Logger
}

func NewLogoutUsecase(
	userRepo user.Repository,
	appRepo application.Repository,
	oauthAccRepo oauthaccesstoken.Repository,
	oauthRefRepo oauthrefreshtoken.Repository,
	logger service.Logger,
) LogoutUsecase {
	return &logoutUsecase{
		userRepo:     userRepo,
		appRepo:      appRepo,
		oauthAccRepo: oauthAccRepo,
		oauthRefRepo: oauthRefRepo,
		logger:       logger,
	}
}

func (uc *logoutUsecase) Execute(ctx context.Context, params LogoutParams) error {
	u, err := uc.userRepo.GetByID(ctx, params.UserID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query user",
			"action", "LOGOUT",
			"user_id", params.UserID,
			"error", err.Error(),
		)
		return apperr.NewDirectError(enum.SERVER_ERROR, "failed to query user")
	}

	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query applicaton",
			"action", "LOGOUT",
			"client_id", &params.ClientID,
			"error", err.Error(),
		)
		return apperr.NewDirectError(enum.SERVER_ERROR, "failed to query application")
	}

	revoked, err := uc.revokeAllToken(ctx, u.ID, app.ID)
	if err != nil {
		return err
	}
	if revoked {
		return nil
	}

	return nil
}

func (uc *logoutUsecase) revokeAllToken(ctx context.Context, userID, appID uuid.UUID) (bool, error) {
	if err := uc.oauthAccRepo.RevokeAllByUserAndApplication(ctx, userID, appID); err != nil {
		uc.logger.Error(ctx, "failed to revoke access token",
			"action", "REVOKE", "application_id", appID, "error", err.Error())
		return false, apperr.NewDirectError(enum.SERVER_ERROR, "failed to revoke access token")
	}

	if err := uc.oauthRefRepo.RevokeAllByUserAndApplication(ctx, userID, appID); err != nil {
		uc.logger.Error(ctx, "failed to revoke refresh token",
			"action", "REVOKE", "application_id", appID, "error", err.Error())
		return false, apperr.NewDirectError(enum.SERVER_ERROR, "failed to revoke refresh token")
	}

	return true, nil
}
