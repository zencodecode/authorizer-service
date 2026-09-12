package auth

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/authorizesession"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organization"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationapplication"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationuser"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type (
	ConsentParams struct {
		OrganizationID uuid.UUID
		ChallengeID    string
	}

	ConsentResult struct {
		User        *entity.User
		ChallengeID string
		RedirectURL *string
	}
)

type consentUsecase struct {
	userRepo    user.Repository
	orgUserRepo organizationuser.Repository
	orgAppRepo  organizationapplication.Repository
	orgRepo     organization.Repository
	appRepo     application.Repository
	sessionRepo authorizesession.Repository
	authCodeSvc service.AuthorizationCode
	logger      service.Logger
}

func NewConsentUsecase(
	userRepo user.Repository,
	orgUserRepo organizationuser.Repository,
	orgAppRepo organizationapplication.Repository,
	orgRepo organization.Repository,
	appRepo application.Repository,
	sessionRepo authorizesession.Repository,
	authCodeSvc service.AuthorizationCode,
	logger service.Logger,
) ConsentUsecase {
	return &consentUsecase{
		userRepo:    userRepo,
		orgUserRepo: orgUserRepo,
		orgAppRepo:  orgAppRepo,
		orgRepo:     orgRepo,
		appRepo:     appRepo,
		sessionRepo: sessionRepo,
		authCodeSvc: authCodeSvc,
		logger:      logger,
	}
}

func (uc *consentUsecase) Execute(ctx context.Context, params ConsentParams) (*ConsentResult, error) {
	sess, err := uc.sessionRepo.Get(ctx, params.ChallengeID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query session params",
			"action", "CONSENT",
			"challenge_id", params.ChallengeID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query session params")
	}

	u, err := uc.userRepo.GetByID(ctx, *sess.UserID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query user",
			"action", "CONSENT",
			"user_id", *sess.UserID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query user")
	}

	orgUser, err := uc.orgUserRepo.GetByOrganizationAndUser(ctx, params.OrganizationID, u.ID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization user",
			"action", "CONSENT",
			"user_id", u.ID,
			"organization_id", params.OrganizationID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query organization user")
	}

	app, err := uc.appRepo.GetByClientID(ctx, sess.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query applicaton",
			"action", "CONSENT",
			"client_id", &sess.ClientID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query application")
	}

	orgApp, err := uc.orgAppRepo.GetByOrganizationAndApplication(ctx, orgUser.OrganizationID, app.ID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization application",
			"action", "CONSENT",
			"organization_id", params.OrganizationID,
			"application_id", sess.ClientID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query organization application")
	}

	if !orgApp.IsActive {
		return nil, apperr.NewDirectError(enum.ACCESS_DENIED, "organization application is not active")
	}

	org, err := uc.orgRepo.GetByID(ctx, orgApp.OrganizationID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization",
			"action", "CONSENT",
			"organization_id", params.OrganizationID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query organization")
	}

	p := service.IssueAuthorizationCodeParams{
		Session:     sess,
		ChallengeID: params.ChallengeID,
		UserID:      u.ID,
		AppID:       app.ID,
		OrgID:       &org.ID,
	}

	issued, err := uc.authCodeSvc.Issue(ctx, p)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate authorization code",
			"action", "LOGIN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate authorization code")
	}

	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", sess.RedirectURI, issued.Code, sess.State)
	return &ConsentResult{
		User:        u,
		ChallengeID: sess.CodeChallenge,
		RedirectURL: &redirectURL,
	}, nil
}
