package auth

import "github.com/zencodecode/authorizer-service/internal/domain/entity"

type (
	RevokeParams struct {
		TokenType    string
		ClientID     string
		ClientSecret string
	}
	RevokeResult struct {
		User         *entity.User
		AccessToken  string
		RefreshToken string
		Scope        string
		IDToken      string
	}
)
