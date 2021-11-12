package external

import (
	"context"
	"errors"

	"github.com/knadh/koanf"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleBaseURL = "https://www.googleapis.com/oauth2/v3/userinfo"

type googleProvider struct {
	*oauth2.Config
}

type googleUser struct {
	Name          string `json:"name"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	AvatarURL     string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

// NewGoogleProvider creates a Google account provider.
func NewGoogleProvider(konf *koanf.Koanf) (OAuthProvider, error) {

	clientID := konf.String("google.client_id")
	secret := konf.String("google.secret")
	redirectURI := konf.String("google.redirect_url")

	return &googleProvider{
		&oauth2.Config{
			ClientID:     clientID,
			ClientSecret: secret,
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"email",
				"profile",
			},
			RedirectURL: redirectURI,
		},
	}, nil
}

func (g googleProvider) GetOAuthToken(code string) (*oauth2.Token, error) {
	return g.Exchange(oauth2.NoContext, code)
}

func (g googleProvider) GetUserData(ctx context.Context, tok *oauth2.Token) (*UserProvidedData, error) {
	var u googleUser
	if err := makeRequest(ctx, tok, g.Config, googleBaseURL, &u); err != nil {
		return nil, err
	}
	data := &UserProvidedData{
		Email:     u.Email,
		Verified:  u.EmailVerified,
		Name:      u.Name,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		AvatarURL: u.AvatarURL,
		Provider:  "google",
	}

	if data.Email == "" {
		return nil, errors.New("Unable to find email with Google provider")
	}

	return data, nil
}
