package auth_api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/henki/pkg/external"
)

func (a *AuthApi) ExternalProviderRedirect(w http.ResponseWriter, r *http.Request) {
	providerType := r.URL.Query().Get("provider")

	strRedirectURL, err := external.CreateRedirectURL(a.svc.Konf, providerType)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}
	http.Redirect(w, r, strRedirectURL, http.StatusFound)
}

func (a *AuthApi) ExternalProviderCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rq := r.URL.Query()

	extError := rq.Get("error")
	if extError != "" {
		http_utils.SendErrorResponse(errors.New(extError), w, r)
		return
	}

	state := rq.Get("state")
	oauthCode := rq.Get("code")

	userInfo, err := external.AuthenticateUser(ctx, a.svc.Konf, state, oauthCode)

	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	accessToken, err := a.usvc.AuthenticateExternalUser(ctx, userInfo)

	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	bjson, err := json.Marshal(accessToken)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	token64 := base64.StdEncoding.EncodeToString(bjson)

	q := url.Values{}
	q.Set("ticket", token64)

	strURL := a.svc.Konf.String("client.external_login_complete_url")

	strURL += "?" + q.Encode()
	http.Redirect(w, r, strURL, http.StatusFound)
}
