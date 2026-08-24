package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type IssueAuthorizationCodeParams struct {
	Session     *entity.AuthorizeSession
	ChallengeID string
	UserID      uuid.UUID
	AppID       uuid.UUID
	OrgID       *uuid.UUID
}

type IssuedAuthorizationCode struct {
	Code   string
	Record *entity.OAuthAuthorizationCode
}

type AuthorizationCode interface {
	Issue(ctx context.Context, params IssueAuthorizationCodeParams) (*IssuedAuthorizationCode, error)
}
