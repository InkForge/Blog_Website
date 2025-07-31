package domain

import (
	"context"
)

type OAuth2User struct {
	ID             string
	Email          string
	VerifiedEmail  bool
	Name           string
	FirstName      string
	LastName       string
	Picture        string
	Provider       string
	RawData        map[string]interface{} 
}

type OAuth2ProviderConfig struct {
	ClientID       string
	ClientSecret   string
	RedirectURL    string
	Scopes         []string
}

// OAuth2 providers interface
type IOAuth2Provider interface {
	
	Name() string   // provider name
	Authenticate(ctx context.Context, code string) (*OAuth2User, error)
	GetAuthorizationURL(state string) string
}