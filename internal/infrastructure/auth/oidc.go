package auth

import "errors"

type OIDCService interface {
	GetDiscovery(issuerURL string) (*OpenIDConfiguration, error)
}

type OpenIDConfiguration struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	JwksURI                           string   `json:"jwks_uri"`
	RevocationEndpoint                string   `json:"revocation_endpoint"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
}

type oidcService struct{}

func NewOIDCService() OIDCService {
	return &oidcService{}
}

func (s *oidcService) GetDiscovery(issuerURL string) (*OpenIDConfiguration, error) {

	if issuerURL == "" {
		return nil, errors.New("issuerURL is invalid")
	}
	data := &OpenIDConfiguration{
		Issuer:                            issuerURL,
		AuthorizationEndpoint:             issuerURL + "/authorize",
		TokenEndpoint:                     issuerURL + "/token",
		UserinfoEndpoint:                  issuerURL + "/userinfo",
		JwksURI:                           issuerURL + "/.well-known/jwks.json",
		RevocationEndpoint:                issuerURL + "/revoke",
		ScopesSupported:                   []string{"openid", "profile", "email", "offline_access"},
		ResponseTypesSupported:            []string{},
		GrantTypesSupported:               []string{},
		TokenEndpointAuthMethodsSupported: []string{},
		CodeChallengeMethodsSupported:     []string{},
		SubjectTypesSupported:             []string{},
		IDTokenSigningAlgValuesSupported:  []string{},
	}

	return data, nil
}
