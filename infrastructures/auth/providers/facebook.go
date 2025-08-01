package providers

// imports
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/InkForge/Blog_Website/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type facebookProvider struct {
	config *oauth2.Config
}

// creates a new Facebook OAuth2 provider
func NewFacebookProvider(cfg domain.OAuth2ProviderConfig) *facebookProvider {
	if len(cfg.Scopes) == 0 {
		cfg.Scopes = []string{"email", "public_profile"}
	}

	return &facebookProvider{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Scopes:       cfg.Scopes,
			Endpoint:     facebook.Endpoint,
		},
	}
}

func (fbprov *facebookProvider) Name() string {
	return "facebook"
}

func (fbprov *facebookProvider) GetAuthorizationURL(state string) string {
	return fbprov.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (fbprov *facebookProvider) Authenticate(ctx context.Context, code string) (*domain.User, error) {
	
	token, err := fbprov.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("facebook: code exchange failed: %w", err)
	}

	client := fbprov.config.Client(ctx, token)
	
	// facebook requires fields to be specified
	profileURL := fmt.Sprintf(
		"https://graph.facebook.com/me?fields=id,name,first_name,last_name,email,picture.type(large)&access_token=%s",
		url.QueryEscape(token.AccessToken),
	)

	resp, err := client.Get(profileURL)
	if err != nil {
		return nil, fmt.Errorf("facebook: failed getting user info: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("facebook: failed reading response body: %w", err)
	}

	var userInfo struct {
		ID        	string 	  `json:"id"`
		Email     	string 	  `json:"email"`
		Name      	*string   `json:"name"`
		FirstName 	*string   `json:"first_name"`
		LastName  	*string   `json:"last_name"`
		Picture   	struct {
			Data 	struct {
				URL *string `json:"url"`
			} `json:"data"`
		} `json:"picture"`
	}

	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("facebook: failed parsing user info: %w", err)
	}

	// parse raw data
	var rawData map[string]interface{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		rawData = make(map[string]interface{})
	}

	return &domain.User{
		UserID:               userInfo.ID,
		Email:                userInfo.Email,
		IsVerified:           true,        // facebook doesn't provide this info
		Name:                 userInfo.Name,
		FirstName:            userInfo.FirstName,
		LastName:      		  userInfo.LastName,
		ProfilePicture:       userInfo.Picture.Data.URL,
		Provider:             fbprov.Name(),
		RawData:              rawData,
	}, nil
}