package external

import (
	"context"
	"encoding/json"
	"golang.org/x/oauth2"
)

type UserProvidedData struct {
	Email     string
	Verified  bool
	Name      string
	FirstName string
	LastName  string
	AvatarURL string
	Provider  string
}

const (
	avatarURLKey = "avatar_url"
	nameKey      = "full_name"
	aliasKey     = "slug"
)

// OAuthProvider specifies additional methods needed for providers using OAuth
type OAuthProvider interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	GetUserData(context.Context, *oauth2.Token) (*UserProvidedData, error)
	GetOAuthToken(string) (*oauth2.Token, error)
}

type AccessProvider interface {
	AuthenticateExternalUser(ctx context.Context,
		email string, meta map[string]string) (string, error)
}

func makeRequest(ctx context.Context, tok *oauth2.Token, g *oauth2.Config, url string, dst interface{}) error {
	client := g.Client(ctx, tok)
	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(dst); err != nil {
		return err
	}

	return nil
}
