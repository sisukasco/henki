package auth_api

import (
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/henki/pkg/usersvc"
	"log"
	"net/http"
)

func (a *AuthApi) GetToken(w http.ResponseWriter, r *http.Request) {
	http_utils.HandleCall(a.handleGetToken, w, r)
}

func (a *AuthApi) handleGetToken(w http.ResponseWriter, r *http.Request) (interface{}, error) {
	grantType := r.FormValue("grant_type")

	token, err := a.createToken(grantType, w, r)
	if err != nil {
		return nil, err
	}

	cookie := http.Cookie{Name: "simtok",
		Value:  token.Token,
		MaxAge: token.ExpiresIn,
	}
	domain := a.svc.Konf.String("auth.cookie.domain")
	if len(domain) > 0 {
		cookie.Domain = domain
	}

	http.SetCookie(w, &cookie)
	log.Printf("Setting cookie to access token %v", token.Token)

	return token, nil
}

func (a *AuthApi) createToken(grantType string, w http.ResponseWriter, r *http.Request) (*usersvc.AccessTokenResponse, error) {
	switch grantType {
	case "password":
		return a.PasswordLogin(w, r)
	case "refresh_token":
		return a.RenewRefreshToken(w, r)
	default:
		return nil, http_utils.OauthError("unsupported_grant_type", "")
	}
}

func (a *AuthApi) PasswordLogin(w http.ResponseWriter, r *http.Request) (*usersvc.AccessTokenResponse, error) {

	var params usersvc.PasswordLoginParams
	params.Username = r.FormValue("username")
	params.Password = r.FormValue("password")

	return a.usvc.PasswordLogin(r.Context(), &params)
}

func (a *AuthApi) RenewRefreshToken(w http.ResponseWriter, r *http.Request) (*usersvc.AccessTokenResponse, error) {
	tokenStr := r.FormValue("refresh_token")

	return a.usvc.RenewRefreshToken(r.Context(), tokenStr)
}
