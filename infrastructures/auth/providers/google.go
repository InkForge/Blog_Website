package providers

// imports
import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/InkForge/Blog_Website/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleProvider struct {
	config *oauth2.Config
}

// creates a new Google OAuth2 provider
func NewGoogleProvider(confg domain.OAuth2ProviderConfig) *googleProvider {
	if len(confg.Scopes) == 0 {
		confg.Scopes = []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		}
	}

	return &googleProvider{
		config: &oauth2.Config{
			ClientID:     confg.ClientID,
			ClientSecret: confg.ClientSecret,
			RedirectURL:  confg.RedirectURL,
			Scopes:       confg.Scopes,
			Endpoint:     google.Endpoint,
		},
	}
}

func (ggprov *googleProvider) Name() string {
	return "google"
}

func (ggprov *googleProvider) GetAuthorizationURL(state string) string {
	return ggprov.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (ggprov *googleProvider) Authenticate(ctx context.Context, code string) (*domain.OAuth2User, error) {
	
	token, err := ggprov.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("google: code exchange failed: %w", err)
	}

	client := ggprov.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("google: failed getting user info: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("google: failed reading response body: %w", err)
	}

	var userInfo struct {
		ID            	string    `json:"id"`
		Email         	string    `json:"email"`
		VerifiedEmail 	bool      `json:"verified_email"`
		Name          	string 	  `json:"name"`
		GivenName     	string 	  `json:"given_name"`
		FamilyName    	string 	  `json:"family_name"`
		Picture       	string 	  `json:"picture"`
	}

	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("google: failed parsing user info: %w", err)
	}

	// parse raw data
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		rawData = make(map[string]interface{})
	}

	return &domain.OAuth2User{
		ID:            userInfo.ID,
		Email:         userInfo.Email,
		VerifiedEmail: userInfo.VerifiedEmail,
		Name:          userInfo.Name,
		FirstName:     userInfo.GivenName,
		LastName:      userInfo.FamilyName,
		Picture:       userInfo.Picture,
		Provider:      ggprov.Name(),
		RawData:       rawData,
	}, nil
}