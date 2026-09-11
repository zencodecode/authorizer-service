package auth

// AuthorizeResponse represents the authorization response (swagger doc only).
// The actual endpoint issues a 302 redirect.
type AuthorizeResponse struct {
	LoginChallenge string `json:"login_challenge"`
}

// ErrorResponse represents an error response (swagger doc only).
type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}
