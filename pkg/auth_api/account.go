package auth_api

import (
	"encoding/json"
	"errors"
	"github.com/sisukasco/commons/crypto"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/commons/utils"
	"github.com/sisukasco/henki/pkg/usersvc"
	"log"
	"net/http"
	"time"
)

func (a *AuthApi) UpdateAccount(w http.ResponseWriter, r *http.Request) {

	http_utils.HandleCall(a.handleUpdateAccount, w, r)
}

func (a *AuthApi) handleUpdateAccount(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {
	ua := &usersvc.AppUser{}

	jsonDecoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := jsonDecoder.Decode(ua)

	if err != nil {
		return nil, err
	}
	log.Printf("Auth Received UserAccountUpdate request.\n%v\n", utils.ToJSONString(ua))

	signx, err := http_utils.ExtractSignature(r)
	if err != nil {
		return nil, http_utils.ForbiddenError("Signature is required for this request")
	}

	key := a.svc.Konf.String("accountx.api_key")
	if len(key) < 5 {
		return nil, errors.New("Account API Key is not set. Path: accountx.api_key")
	}

	err = crypto.VerifySignature(ua, key, signx)
	if err != nil {
		return nil, http_utils.ForbiddenError("Bad signature for this request")
	}

	resp, err := a.usvc.UpdateUserAccount(r.Context(), ua.AccountInfo.Contact.Email, ua)
	if err != nil {
		log.Printf("Auth: Error updating user account %v", err)
		return nil, err
	}

	return resp, nil
}

type RemoveAccountRequest struct {
	Email     string `json:"email"`
	CreatedAt time.Time
}

func (a *AuthApi) RemoveAccount(w http.ResponseWriter, r *http.Request) {

	http_utils.HandleCall(a.handleRemoveAccount, w, r)
}

func (a *AuthApi) handleRemoveAccount(w http.ResponseWriter,
	r *http.Request) (interface{}, error) {
	rreq := &RemoveAccountRequest{}

	jsonDecoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	err := jsonDecoder.Decode(rreq)

	if err != nil {
		return nil, err
	}
	log.Printf("Auth Received RemoveAccount request.\n%v\n", utils.ToJSONString(rreq))

	signx, err := http_utils.ExtractSignature(r)
	if err != nil {
		return nil, http_utils.ForbiddenError("Signature is required for this request")
	}

	key := a.svc.Konf.String("accountx.api_key")
	if len(key) < 5 {
		return nil, errors.New("Account API Key is not set")
	}

	err = crypto.VerifySignature(rreq, key, signx)
	if err != nil {
		return nil, http_utils.ForbiddenError("Bad signature for this request")
	}

	resp, err := a.usvc.SwitchToFreeAccount(r.Context(), rreq.Email)
	if err != nil {
		log.Printf("Auth: Error switching user account %v", err)
		return nil, err
	}

	return resp, nil
}
