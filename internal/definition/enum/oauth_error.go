package enum

type OAuthError string

const (
	INVALID_REQUEST           OAuthError = "invalid_request"
	UNAUTHORIZED_CLIENT       OAuthError = "unauthorized_client"
	ACCESS_DENIED             OAuthError = "access_denied"
	UNSUPPORTED_RESPONSE_TYPE OAuthError = "unsupported_response_type"
	INVALID_SCOPE             OAuthError = "invalid_scope"
	SERVER_ERROR              OAuthError = "server_error"
	TEMPORARILY_UNAVAILABLE   OAuthError = "temporarily_unavailable"
)
