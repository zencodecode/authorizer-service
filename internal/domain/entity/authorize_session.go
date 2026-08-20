package entity

type AuthorizeSession struct {
	ClientID            string
	RedirectURI         string
	Scope               []string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
}
