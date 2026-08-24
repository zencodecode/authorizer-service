package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/authorizesession"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthauthorizationcode"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
	"github.com/zencodecode/authorizer-service/pkg/randutil"
)

type authorizationCode struct {
	oauthCodeRepo oauthauthorizationcode.Repository
	sessionRepo   authorizesession.Repository
}

func NewAuthorizationCode(
	oauthCodeRepo oauthauthorizationcode.Repository,
	sessionRepo authorizesession.Repository,
) service.AuthorizationCode {
	return &authorizationCode{
		oauthCodeRepo: oauthCodeRepo,
		sessionRepo:   sessionRepo,
	}
}

func (s *authorizationCode) Issue(ctx context.Context, params service.IssueAuthorizationCodeParams) (*service.IssuedAuthorizationCode, error) {
	code, err := randutil.GenerateRandomString(32)
	if err != nil {
		return nil, errors.New("failed to generate codebytes")
	}

	codeHash := hash.HashSHA256(code)

	authCode := &entity.OAuthAuthorizationCode{
		ID:                  uuid.Must(uuid.NewV7()),
		CodeHash:            codeHash,
		UserID:              params.UserID,
		OrganizationID:      params.OrgID,
		ApplicationID:       params.AppID,
		RedirectURI:         params.Session.RedirectURI,
		CodeChallenge:       params.Session.CodeChallenge,
		CodeChallengeMethod: "S256",
		Scopes:              params.Session.Scope,
		ExpiresAt:           time.Now().Add(60 * time.Second),
	}

	if err := s.oauthCodeRepo.Create(ctx, authCode); err != nil {
		return nil, errors.New("failed to persist oauth authorize code")
	}

	_ = s.sessionRepo.Delete(ctx, params.ChallengeID)
	return &service.IssuedAuthorizationCode{
		Code:   code,
		Record: authCode,
	}, nil
}
