package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/k0kubun/pp"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/authorizesession"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthauthorizationcode"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organization"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationapplication"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationuser"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type (
	ConsentParams struct {
		OrgID       uuid.UUID
		ChallengeID string
	}
	ConsentResult struct {
		OrgID       string
		ChallengeID string
	}
)

type consentUsecase struct {
	userRepo      user.Repository
	orgUserRepo   organizationuser.Repository
	orgRepo       organization.Repository
	orgAppRepo    organizationapplication.Repository
	oauthCodeRepo oauthauthorizationcode.Repository
	sessionRepo   authorizesession.Repository
	logger        service.Logger
}

func NewConsentUsecase(
	userRepo user.Repository,
	orgUserRepo organizationuser.Repository,
	orgRepo organization.Repository,
	orgAppRepo organizationapplication.Repository,
	oauthCodeRepo oauthauthorizationcode.Repository,
	sessionRepo authorizesession.Repository,
	logger service.Logger,
) ConsentUsecase {
	return &consentUsecase{
		userRepo:      userRepo,
		orgUserRepo:   orgUserRepo,
		orgRepo:       orgRepo,
		orgAppRepo:    orgAppRepo,
		oauthCodeRepo: oauthCodeRepo,
		sessionRepo:   sessionRepo,
		logger:        logger,
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
		return nil, newAuthError("server_error", "failed to query session params")
	}

	user, err := uc.userRepo.GetByID(ctx, *sess.UserID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query user",
			"action", "CONSENT",
			"user_id", *sess.UserID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query user")
	}

	orgUser, err := uc.orgUserRepo.GetByOrganizationAndUser(ctx, params.OrgID, user.ID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization user",
			"action", "CONSENT",
			"user_id", user.ID,
			"organization_id", params.OrgID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query organization user")
	}

	appID, err := uuid.Parse(sess.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to parse UUID",
			"action", "CONSENT",
			"application_id", sess.ClientID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to parse UUID")
	}

	orgApp, err := uc.orgAppRepo.GetByOrganizationAndApplication(ctx, orgUser.OrganizationID, appID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization application",
			"action", "CONSENT",
			"organization_id", params.OrgID,
			"application_id", sess.ClientID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query organization application")
	}

	pp.Println(orgApp)

	return nil, nil
}
