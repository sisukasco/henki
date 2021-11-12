package auth_api

import (
	"encoding/json"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/henki/pkg/usersvc"
	"net/http"
)

func (a *AuthApi) Signup(w http.ResponseWriter, r *http.Request) {
	http_utils.HandleCall(a.HandleSignup, w, r)
}

func (a *AuthApi) HandleSignup(w http.ResponseWriter, r *http.Request) (interface{}, error) {

	params := &usersvc.SignupParams{}
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(params)
	if err != nil {
		return nil, err
	}

	_, err = a.usvc.SignupNewUser(r.Context(), params)
	if err != nil {
		return nil, err
	}

	return &http_utils.StatusResponse{"ok"}, nil
}
