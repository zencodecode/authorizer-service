package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/authorizesession"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthauthorizationcode"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organization"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationapplication"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationuser"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
	"github.com/zencodecode/authorizer-service/pkg/randutil"
)

type (
	ConsentParams struct {
		OrgID       uuid.UUID
		ChallengeID string
	}
	ConsentResult struct {
		User              *entity.User
		Organizations     []*entity.Organization
		ChallengeID       string
		AuthorizationCode *string
		ExpiresAt         *int64
		RedirectURL       *string
	}
)

type consentUsecase struct {
	userRepo      user.Repository
	orgUserRepo   organizationuser.Repository
	orgAppRepo    organizationapplication.Repository
	orgRepo       organization.Repository
	appRepo       application.Repository
	oauthCodeRepo oauthauthorizationcode.Repository
	sessionRepo   authorizesession.Repository
	logger        service.Logger
}

func NewConsentUsecase(
	userRepo user.Repository,
	orgUserRepo organizationuser.Repository,
	orgAppRepo organizationapplication.Repository,
	orgRepo organization.Repository,
	appRepo application.Repository,
	oauthCodeRepo oauthauthorizationcode.Repository,
	sessionRepo authorizesession.Repository,
	logger service.Logger,
) ConsentUsecase {
	return &consentUsecase{
		userRepo:      userRepo,
		orgUserRepo:   orgUserRepo,
		orgAppRepo:    orgAppRepo,
		orgRepo:       orgRepo,
		appRepo:       appRepo,
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

	u, err := uc.userRepo.GetByID(ctx, *sess.UserID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query user",
			"action", "CONSENT",
			"user_id", *sess.UserID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query user")
	}

	orgUser, err := uc.orgUserRepo.GetByOrganizationAndUser(ctx, params.OrgID, u.ID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization user",
			"action", "CONSENT",
			"user_id", u.ID,
			"organization_id", params.OrgID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query organization user")
	}

	app, err := uc.appRepo.GetByClientID(ctx, sess.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query applicaton",
			"action", "CONSENT",
			"client_id", &sess.ClientID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query applicaton")
	}

	orgApp, err := uc.orgAppRepo.GetByOrganizationAndApplication(ctx, orgUser.OrganizationID, app.ID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization application",
			"action", "CONSENT",
			"organization_id", params.OrgID,
			"application_id", sess.ClientID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query organization application")
	}

	if !orgApp.IsActive {
		return nil, newAuthError("access_denied", "organization application is not active")
	}

	org, err := uc.orgRepo.GetByID(ctx, orgApp.OrganizationID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization",
			"action", "CONSENT",
			"organization_id", params.OrgID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query organization")
	}

	code, err := uc.generateAuthorizationCode(ctx, sess, u.ID, app.ID, &org.ID)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate authorization code",
			"action", "LOGIN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to generate authorization code")
	}

	exp := code.ExpiresAt.Unix()
	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", sess.RedirectURI, sess.CodeChallenge, sess.State)
	return &ConsentResult{
		User:              u,
		AuthorizationCode: &code.CodeHash,
		ExpiresAt:         &exp,
		ChallengeID:       sess.CodeChallenge,
		RedirectURL:       &redirectURL,
	}, nil
}

func (uc *consentUsecase) generateAuthorizationCode(
	ctx context.Context,
	sess *entity.AuthorizeSession,
	userID uuid.UUID,
	appID uuid.UUID,
	orgID *uuid.UUID,
) (*entity.OAuthAuthorizationCode, error) {
	code, err := randutil.GenerateRandomString(32)
	if err != nil {
		return nil, newAuthError("server_error", "failed to generate codebytes")
	}

	codeHash := hash.HashSHA256(code)

	authCode := &entity.OAuthAuthorizationCode{
		ID:                  uuid.Must(uuid.NewV7()),
		CodeHash:            codeHash,
		UserID:              userID,
		OrganizationID:      orgID,
		ApplicationID:       appID,
		RedirectURI:         sess.RedirectURI,
		CodeChallenge:       sess.CodeChallenge,
		CodeChallengeMethod: "S256",
		Scopes:              sess.Scope,
		ExpiresAt:           time.Now().Add(60 * time.Second),
	}

	if err := uc.oauthCodeRepo.Create(ctx, authCode); err != nil {
		uc.logger.Error(ctx, "failed to save authorize session",
			"action", "CONSENT",
			"error", err.Error())
		return nil, newAuthError("server_error", "failed to persist oauth authorize code")
	}

	_ = uc.sessionRepo.Delete(ctx, sess.CodeChallenge)

	return authCode, nil
}
