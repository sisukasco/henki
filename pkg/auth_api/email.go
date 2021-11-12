package auth_api

import (
	"encoding/json"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/henki/pkg/auth_utils"
	"log"
	"net/http"
	"strings"
)

type ConfirmEmailParams struct {
	Code string `json:"code"`
}

func (a *AuthApi) ConfirmEmail(w http.ResponseWriter, r *http.Request) {

	params := &ConfirmEmailParams{}

	jsonDecoder := json.NewDecoder(r.Body)

	err := jsonDecoder.Decode(params)

	if err != nil {
		err = http_utils.InternalServerError("Error decoding request").WithInternalError(err)
		http_utils.SendErrorResponse(err, w, r)
	}

	code := strings.TrimSpace(params.Code)

	resp, err := a.usvc.ConfirmUserEmail(r.Context(), code)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
		//return http_utils.BadRequestError("Could not read Signup params: %v", err)
	}
	err = http_utils.SendJSON(w, http.StatusOK, resp)
	if err != nil {
		log.Printf("Error sending response %v", err)
	}
}

type UpdateEmailParams struct {
	Token string `json:"token"`
}

func (a *AuthApi) CommitEmailUpdate(w http.ResponseWriter, r *http.Request) {
	params := &UpdateEmailParams{}

	jsonDecoder := json.NewDecoder(r.Body)

	err := jsonDecoder.Decode(params)

	if err != nil {
		err = http_utils.InternalServerError("Error decoding request").WithInternalError(err)
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	err = a.usvc.CompleteEmailUpdate(r.Context(), params.Token)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}
	err = http_utils.SendJSON(w, http.StatusOK, &http_utils.StatusResponse{"ok"})
	if err != nil {
		log.Printf("Error sending response %v", err)
	}
}

func (a *AuthApi) ResendConfirmEmail(w http.ResponseWriter, r *http.Request) {
	user, err := auth_utils.GetUserFromRequest(a.svc, r)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	err = a.usvc.SendEmailConfirmationRequest(r.Context(), user.ID)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
	}
	err = http_utils.SendJSON(w, http.StatusOK, &http_utils.StatusResponse{"ok"})
	if err != nil {
		log.Printf("Error sending response %v", err)
	}

}
