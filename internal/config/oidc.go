package config

import (
	"fmt"
	"time"

	"github.com/zencodecode/authorizer-service/pkg/envutil"
)

type OIDC struct {
	JWT                     JWT
	Issuer                  string
	AuthorizationCodeExpiry time.Duration
	RequirePKCE             bool
}

func LoadOIDCConfig(jwtCfg JWT) (OIDC, error) {
	issuer := envutil.Get("OAUTH_ISSUER_URL", "")
	if issuer == "" {
		return OIDC{}, fmt.Errorf("OAUTH_ISSUER_URL must be set")
	}

	codeExpiry := envutil.GetDuration("OAUTH_CODE_EXPIRY", 60*time.Second)
	requirePKCE := envutil.GetAsBool("OAUTH_REQUIRE_PKCE", true)

	return OIDC{
		JWT:                     jwtCfg,
		Issuer:                  issuer,
		AuthorizationCodeExpiry: codeExpiry,
		RequirePKCE:             requirePKCE,
	}, nil
}
