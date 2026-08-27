package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/auditlog"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/authorizesession"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthauthorizationcode"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organization"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organizationuser"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type NextAction string

const CONSENT NextAction = "consent"
const REDIRECT NextAction = "redirect"

type (
	LoginParams struct {
		Email       string
		Password    string
		ChallengeID string
	}

	LoginResult struct {
		User          *entity.User
		Organizations []*entity.Organization
		ChallengeID   string
		RedirectURL   *string
		NextStep      NextAction
	}
)

type loginUsecase struct {
	userRepo      user.Repository
	sessionRepo   authorizesession.Repository
	appRepo       application.Repository
	orgUserRepo   organizationuser.Repository
	orgRepo       organization.Repository
	oauthCodeRepo oauthauthorizationcode.Repository
	auditLogRepo  auditlog.Repository
	authCodeSvc   service.AuthorizationCode
	logger        service.Logger
}

func NewLoginUsecase(
	userRepo user.Repository,
	sessionRepo authorizesession.Repository,
	appRepo application.Repository,
	orgUserRepo organizationuser.Repository,
	orgRepo organization.Repository,
	oauthCodeRepo oauthauthorizationcode.Repository,
	auditLogRepo auditlog.Repository,
	authCodeSvc service.AuthorizationCode,
	logger service.Logger,
) LoginUsecase {
	return &loginUsecase{
		userRepo:      userRepo,
		sessionRepo:   sessionRepo,
		appRepo:       appRepo,
		orgUserRepo:   orgUserRepo,
		orgRepo:       orgRepo,
		oauthCodeRepo: oauthCodeRepo,
		auditLogRepo:  auditLogRepo,
		authCodeSvc:   authCodeSvc,
		logger:        logger,
	}
}

func (uc *loginUsecase) Execute(ctx context.Context, params LoginParams) (*LoginResult, error) {

	sess, err := uc.sessionRepo.Get(ctx, params.ChallengeID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query session params",
			"action", "LOGIN",
			"challenge_id", params.ChallengeID,
			"error", err.Error(),
		)
		return nil, apperr.NewFatalError(enum.SERVER_ERROR, "failed to query session params")
	}

	u, err := uc.userRepo.GetByEmail(ctx, params.Email)
	if err != nil {
		uc.logger.Warn(ctx, "failed to query user by email",
			"action", "LOGIN",
			"email", params.Email,
			"error", err.Error(),
		)
		return nil, apperr.NewRedirectableError(enum.INVALID_GRANT,
			"email or password is invalid", sess.RedirectURI, sess.State)
	}

	if !hash.CheckHash(u.PasswordHash, params.Password) {
		uc.logger.Warn(ctx, "invalid password",
			"action", "LOGIN",
			"user_id", u.ID,
		)
		return nil, apperr.NewRedirectableError(enum.INVALID_GRANT,
			"email or password is invalid", sess.RedirectURI, sess.State)
	}

	if u.Status != "active" {
		uc.logger.Warn(ctx, "login attempt on inactive account",
			"action", "LOGIN",
			"user_id", u.ID,
			"status", u.Status,
		)
		return nil, apperr.NewFatalError(enum.ACCESS_DENIED, "account is suspended")
	}

	if u.EmailVerifiedAt == nil {
		uc.logger.Warn(ctx, "login attempt on unverified account",
			"action", "LOGIN",
			"user_id", u.ID,
			"status", u.Status,
		)
		return nil, apperr.NewFatalError(enum.ACCESS_DENIED, "email is not verified")
	}

	app, err := uc.appRepo.GetByClientID(ctx, sess.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query application",
			"action", "LOGIN",
			"client_id", sess.ClientID,
			"error", err.Error(),
		)
		return nil, apperr.NewFatalError(enum.SERVER_ERROR, "failed to query application")
	}

	sess.UserID = &u.ID
	if err := uc.sessionRepo.Save(ctx, params.ChallengeID, *sess, 10*time.Minute); err != nil {
		return nil, apperr.NewFatalError(enum.SERVER_ERROR, "failed to persist session")
	}

	if app.RequiresOrganization {
		var orgsUser []*entity.OrganizationUser
		orgsUser, err = uc.orgUserRepo.ListOrganizationsByUser(ctx, u.ID)
		if err != nil {
			uc.logger.Error(ctx, "failed to query organization user",
				"action", "LOGIN",
				"user_id", u.ID,
				"error", err.Error(),
			)
			return nil, apperr.NewFatalError(enum.SERVER_ERROR, "failed to query organization user")
		}
		orgSet := make([]*entity.Organization, 0, len(orgsUser))

		for _, ou := range orgsUser {
			org, err := uc.orgRepo.GetByID(ctx, ou.OrganizationID)
			if err != nil {
				uc.logger.Error(ctx, "failed to query organization",
					"action", "LOGIN",
					"organization_id", ou.OrganizationID,
					"error", err.Error(),
				)
				return nil, apperr.NewFatalError(enum.SERVER_ERROR, "failed to query organization")
			}

			orgSet = append(orgSet, org)
		}
		return &LoginResult{
			User:          u,
			Organizations: orgSet,
			ChallengeID:   params.ChallengeID,
			NextStep:      CONSENT,
		}, nil

	}

	p := service.IssueAuthorizationCodeParams{
		Session:     sess,
		ChallengeID: params.ChallengeID,
		UserID:      u.ID,
		AppID:       app.ID,
		OrgID:       nil,
	}

	issued, err := uc.authCodeSvc.Issue(ctx, p)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate authorization code",
			"action", "LOGIN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, apperr.NewFatalError(enum.SERVER_ERROR, "failed to generate authorization code")
	}

	redirectURL := fmt.Sprintf("%s?code=%s&state=%s", sess.RedirectURI, issued.Code, sess.State)
	return &LoginResult{
		User:        u,
		ChallengeID: params.ChallengeID,
		RedirectURL: &redirectURL,
		NextStep:    REDIRECT,
	}, nil
}
