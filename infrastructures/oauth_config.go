package infrastructures

import (
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/InkForge/Blog_Website/domain"
)

func BuildProviderConfigs() (map[string]domain.OAuth2ProviderConfig, error) {
	configs, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return map[string]domain.OAuth2ProviderConfig{
		"google": {
			ClientID:     configs.GoogleClientID,
			ClientSecret: configs.GoogleClientSecret,
			RedirectURL:  configs.BaseURL + configs.GoogleRedirectURL,
			Scopes:       []string{"profile", "email"},
			Endpoint:     google.Endpoint,
		},
		"github": {
			ClientID:     configs.GithubClientID,
			ClientSecret: configs.GithubClientSecret,
			RedirectURL:  configs.BaseURL + configs.GithubRedirectURL,
			Scopes:       []string{"user:email"},
			Endpoint:     github.Endpoint,
		},
		"facebook": {
			ClientID:     configs.FacebookClientID,
			ClientSecret: configs.FacebookClientSecret,
			RedirectURL:  configs.BaseURL + configs.FacebookRedirectURL,
			Scopes:       []string{"email", "public_profile"},
			Endpoint:     facebook.Endpoint,
		},
	}, nil
}
