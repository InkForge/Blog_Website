package providers

// imports
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/InkForge/Blog_Website/domain"
	"golang.org/x/oauth2"
)

type githubProvider struct {
	config *oauth2.Config
}

// creates a new GitHub OAuth2 provider
func NewGitHubProvider(confg domain.OAuth2ProviderConfig) *githubProvider {
	if len(confg.Scopes) == 0 {
		confg.Scopes = []string{"user:email"}
	}

	return &githubProvider{
		config: &oauth2.Config{
			ClientID:     confg.ClientID,
			ClientSecret: confg.ClientSecret,
			RedirectURL:  confg.RedirectURL,
			Scopes:       confg.Scopes,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://github.com/login/oauth/authorize",
				TokenURL: "https://github.com/login/oauth/access_token",
			},
		},
	}
}

func (ghprov *githubProvider) Name() string {
	return "github"
}

func (ghprov *githubProvider) GetAuthorizationURL(state string) string {
	return ghprov.config.AuthCodeURL(state)
}

func (ghprov *githubProvider) Authenticate(ctx context.Context, code string) (*domain.OAuth2User, error) {
	
	token, err := ghprov.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("github: code exchange failed: %w", err)
	}

	client := ghprov.config.Client(ctx, token)
	
	// get user profile
	profileResp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("github: failed getting user profile: %w", err)
	}
	defer profileResp.Body.Close()

	profileData, err := io.ReadAll(profileResp.Body)
	if err != nil {
		return nil, fmt.Errorf("github: failed reading profile response: %w", err)
	}

	var profile struct {
		ID      	int      `json:"id"`
		Login   	string   `json:"login"`
		Name    	string   `json:"name"`
		Email   	string   `json:"email"`
		Picture 	string   `json:"avatar_url"`
	}

	if err := json.Unmarshal(profileData, &profile); err != nil {
		return nil, fmt.Errorf("github: failed parsing profile: %w", err)
	}

	// get primary email if not returned in profile
	if profile.Email == "" {
		emailsResp, err := client.Get("https://api.github.com/user/emails")
		if err != nil {
			return nil, fmt.Errorf("github: failed getting user emails: %w", err)
		}
		defer emailsResp.Body.Close()

		emailsData, err := io.ReadAll(emailsResp.Body)
		if err != nil {
			return nil, fmt.Errorf("github: failed reading emails response: %w", err)
		}

		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}

		if err := json.Unmarshal(emailsData, &emails); err != nil {
			return nil, fmt.Errorf("github: failed parsing emails: %w", err)
		}

		for _, email := range emails {
			if email.Primary {
				profile.Email = email.Email
				break
			}
		}
	}

	// parse raw data
	var rawData map[string]interface{}
	if err := json.Unmarshal(profileData, &rawData); err != nil {
		rawData = make(map[string]interface{})
	}

	return &domain.OAuth2User{
		ID:            fmt.Sprintf("%d", profile.ID),
		Email:         profile.Email,
		VerifiedEmail: true,      // GitHub doesn't provide this info
		Name:          profile.Name,
		FirstName:     "",        // GitHub doesn't provide first name separately
		LastName:      "",		  // GitHub doesn't provide last name separately
		Picture:       profile.Picture,
		Provider:      ghprov.Name(),
		RawData:       rawData,
	}, nil
}