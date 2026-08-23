package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	"github.com/zencodecode/authorizer-service/pkg/randutil"
)

type (
	LoginParams struct {
		Email       string
		Password    string
		ChallengeID string
	}

	LoginResult struct {
		User              *entity.User
		Organizations     []*entity.Organization
		AuthorizationCode *string
		ExpiresAt         *int64
		Session           *entity.AuthorizeSession
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
	// userRoleRepo     userrole.Repository
	// rolePermRepo     rolepermission.Repository
	// refreshTokenRepo oauthrefreshtoken.Repository
	// jwtService       service.JWTService
	logger service.Logger
}

func NewLoginUsecase(
	userRepo user.Repository,
	sessionRepo authorizesession.Repository,
	appRepo application.Repository,
	orgUserRepo organizationuser.Repository,
	orgRepo organization.Repository,
	oauthCodeRepo oauthauthorizationcode.Repository,
	auditLogRepo auditlog.Repository,
	// userRoleRepo userrole.Repository,
	// rolePermRepo rolepermission.Repository,
	// refreshTokenRepo oauthrefreshtoken.Repository,
	// jwtService service.JWTService,
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
		// userRoleRepo:     userRoleRepo,
		// rolePermRepo:     rolePermRepo,
		// refreshTokenRepo: refreshTokenRepo,
		// jwtService:       jwtService,
		logger: logger,
	}
}

func (uc *loginUsecase) Execute(ctx context.Context, params LoginParams) (*LoginResult, error) {
	u, err := uc.userRepo.GetByEmail(ctx, params.Email)
	if err != nil {
		uc.logger.Warn(ctx, "failed to fetch user by email",
			"action", "LOGIN",
			"email", params.Email,
			"error", err.Error(),
		)
		return nil, ErrInvalidCredentials
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}

	if !hash.CheckHash(u.PasswordHash, params.Password) {
		uc.logger.Warn(ctx, "invalid password",
			"action", "LOGIN",
			"user_id", u.ID,
		)
		return nil, ErrInvalidCredentials
	}

	if u.Status != "active" {
		uc.logger.Warn(ctx, "login attempt on inactive account",
			"action", "LOGIN",
			"user_id", u.ID,
			"status", u.Status,
		)
		return nil, ErrAccountSuspended
	}

	if u.EmailVerifiedAt == nil {
		uc.logger.Warn(ctx, "login attempt on unverified account",
			"action", "LOGIN",
			"user_id", u.ID,
			"status", u.Status,
		)
		return nil, ErrEmailNotVerified
	}

	sess, err := uc.sessionRepo.Get(ctx, params.ChallengeID)
	if err != nil {
		uc.logger.Error(ctx, "failed to fetch session params",
			"action", "LOGIN",
			"challenge_id", params.ChallengeID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to fetch session params")
	}

	app, err := uc.appRepo.GetByClientID(ctx, sess.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to fetch application",
			"action", "LOGIN",
			"client_id", sess.ClientID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to fetch application")
	}
	var orgsUser []*entity.OrganizationUser
	if app.RequiresOrganization {
		orgsUser, err = uc.orgUserRepo.ListOrganizationsByUser(ctx, u.ID)
		if err != nil {
			uc.logger.Error(ctx, "failed to fetch organization user",
				"action", "LOGIN",
				"user_id", u.ID,
				"error", err.Error(),
			)
			return nil, newAuthError("server_error", "failed to fetch organization user")
		}

	}

	orgSet := make([]*entity.Organization, 0, len(orgsUser))
	if orgsUser != nil {
		for _, ou := range orgsUser {
			org, err := uc.orgRepo.GetByID(ctx, ou.OrganizationID)
			if err != nil {
				uc.logger.Error(ctx, "failed to fetch organization",
					"action", "LOGIN",
					"organization_id", ou.OrganizationID,
					"error", err.Error(),
				)
				return nil, newAuthError("server_error", "failed to fetch organization")
			}

			orgSet = append(orgSet, org)
		}
		return &LoginResult{
			User:              u,
			Organizations:     orgSet,
			AuthorizationCode: nil,
			ExpiresAt:         nil,
			Session:           sess,
		}, nil

	}

	code, err := uc.generateAuthorizationCode(ctx, sess, u.ID, app.ID, &uuid.Nil)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate authorization code",
			"action", "LOGIN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to generate authorization code")
	}

	// roles, err := uc.userRoleRepo.ListRolesByUser(ctx, u.ID, params.OrgID)
	// if err != nil {
	// 	uc.logger.Error(ctx, "failed to fetch roles",
	// 		"action", "LOGIN",
	// 		"user_id", u.ID,
	// 		"error", err.Error(),
	// 	)
	// 	return nil, newAuthError("server_error", "failed to fetch user roles")
	// }

	// roleSlugs := make([]string, 0, len(roles))
	// permSet := make(map[string]struct{})

	// for _, r := range roles {
	// 	if r.ApplicationID != params.ApplicationID {
	// 		continue
	// 	}
	// 	roleSlugs = append(roleSlugs, r.Slug)

	// 	perms, _ := uc.rolePermRepo.ListPermissionsByRole(ctx, r.ID)
	// 	for _, p := range perms {
	// 		permSet[p.Slug] = struct{}{}
	// 	}
	// }

	// permSlugs := make([]string, 0, len(permSet))
	// for k := range permSet {
	// 	permSlugs = append(permSlugs, k)
	// }

	// now := time.Now()
	// expiresIn := 15 * time.Minute

	// claims := &entity.Claims{
	// 	RegisteredClaims: jwt.RegisteredClaims{
	// 		Issuer:    "authorizer-service",
	// 		Subject:   u.ID.String(),
	// 		Audience:  jwt.ClaimStrings{params.ApplicationID.String()},
	// 		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
	// 		IssuedAt:  jwt.NewNumericDate(now),
	// 		ID:        uuid.Must(uuid.NewV7()).String(),
	// 	},
	// 	Name:        u.Name,
	// 	Email:       u.Email,
	// 	Scopes:      []string{"openid", "profile", "email"},
	// 	Roles:       roleSlugs,
	// 	Permissions: permSlugs,
	// }

	// if params.OrgID != nil {
	// 	orgIDStr := params.OrgID.String()
	// 	claims.OrgID = &orgIDStr
	// }

	// accessToken, err := uc.jwtService.GenerateAccessToken(ctx, claims)
	// if err != nil {
	// 	uc.logger.Error(ctx, "failed to generate access token",
	// 		"action", "LOGIN",
	// 		"user_id", u.ID,
	// 		"error", err.Error(),
	// 	)
	// 	return nil, newAuthError("server_error", "failed to generate access token")
	// }

	// refreshToken, err := uc.jwtService.GenerateRefreshToken()
	// if err != nil {
	// 	uc.logger.Error(ctx, "failed to generate refresh token",
	// 		"action", "LOGIN",
	// 		"user_id", u.ID,
	// 		"error", err.Error(),
	// 	)
	// 	return nil, newAuthError("server_error", "failed to generate refresh token")
	// }

	// // TODO: Store refresh token hash via oauthrefreshtoken repository
	exp := code.ExpiresAt.Unix()
	return &LoginResult{
		User:              u,
		Organizations:     orgSet,
		AuthorizationCode: &code.CodeHash,
		ExpiresAt:         &exp,
		Session:           sess,
	}, nil
}

func (uc *loginUsecase) generateAuthorizationCode(
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
			"action", "AUTHORIZE",
			"error", err.Error())
		return nil, newAuthError("server_error", "failed to persist oauth authorize code")
	}

	// // 5. Redirect
	// redirectURL := fmt.Sprintf("%s?code=%s&state=%s", params.RedirectURI, code, params.State)
	// // return redirect
	return authCode, nil
}
