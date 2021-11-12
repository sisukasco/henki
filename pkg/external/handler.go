package external

import (
	"context"
	"fmt"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/commons/utils"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/knadh/koanf"
)

type ExternalProviderClaims struct {
	jwt.StandardClaims
	SiteURL  string `json:"site_url"`
	Provider string `json:"provider"`

	//Referrer    string `json:"referrer,omitempty"`
}

func CreateRedirectURL(konf *koanf.Koanf, providerType string) (string, error) {
	providerType = utils.CleanupString(providerType)
	if providerType == "" {
		return "", http_utils.BadRequestError("provider name is required ")
	}

	provider, err := getProvider(konf, providerType)
	if err != nil {
		return "", http_utils.BadRequestError("Unsupported provider: %s", providerType).WithInternalError(err)
	}

	siteURL := konf.String("auth.site.url")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ExternalProviderClaims{

		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		SiteURL:  siteURL,
		Provider: providerType,
		//Referrer:    referrer,
	})
	operatorToken := konf.String("auth.operator_token")

	tokenString, err := token.SignedString([]byte(operatorToken))
	if err != nil {
		return "", http_utils.InternalServerError("Error creating state").WithInternalError(err)
	}

	return provider.AuthCodeURL(tokenString), nil

}

func DecodeJwtClaims(ctx context.Context, konf *koanf.Koanf, state string) (*ExternalProviderClaims, error) {

	operatorToken := konf.String("auth.operator_token")

	claims := ExternalProviderClaims{}
	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}

	_, err := p.ParseWithClaims(state, &claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(operatorToken), nil
		})
	if err != nil {
		return nil, http_utils.InternalServerError("Error parsing auth tokens").WithInternalError(err)
	}
	return &claims, nil
}

func AuthenticateUser(ctx context.Context, konf *koanf.Koanf,
	state string, oauthCode string) (*UserProvidedData, error) {

	claims, err := DecodeJwtClaims(ctx, konf, state)
	if err != nil {
		return nil, err
	}
	provider, err := getProvider(konf, claims.Provider)
	if err != nil {
		return nil, http_utils.InternalServerError("Wrong provider string").WithInternalError(err)
	}

	tok, err := provider.GetOAuthToken(oauthCode)
	if err != nil {
		return nil, http_utils.InternalServerError("can't get auth token").WithInternalError(err)
	}
	userProvidedData, err := provider.GetUserData(ctx, tok)

	if err != nil {
		return nil, http_utils.InternalServerError("can't get user data").WithInternalError(err)
	}

	userProvidedData.Provider = claims.Provider

	return userProvidedData, nil

}

func getProvider(konf *koanf.Koanf, name string) (OAuthProvider, error) {
	name = strings.ToLower(name)

	switch name {
	case "google":
		return NewGoogleProvider(konf)
	default:
		return nil, fmt.Errorf("Provider %s could not be found", name)
	}
}
