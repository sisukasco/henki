package auth_api

import (
	"encoding/json"
	"github.com/sisukasco/commons/http_utils"
	"github.com/sisukasco/henki/pkg/auth_utils"
	"log"
	"net/http"
)

type SendResetPasswordParams struct {
	Email string `json:"email"`
}

func (a *AuthApi) SendResetPassword(w http.ResponseWriter, r *http.Request) {
	params := &SendResetPasswordParams{}

	jsonDecoder := json.NewDecoder(r.Body)

	err := jsonDecoder.Decode(params)

	if err != nil {
		err = http_utils.InternalServerError("Error decoding request").WithInternalError(err)
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	err = a.usvc.InitResetPasswordRequest(r.Context(), params.Email)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}
	err = http_utils.SendJSON(w, http.StatusOK, &http_utils.StatusResponse{"ok"})
	if err != nil {
		log.Printf("Error sending response %v", err)
	}
}

type UpdatePasswordParams struct {
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
}

func (a *AuthApi) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	user, err := auth_utils.GetUserFromRequest(a.svc, r)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	params := &UpdatePasswordParams{}

	jsonDecoder := json.NewDecoder(r.Body)

	err = jsonDecoder.Decode(params)

	if err != nil {
		err = http_utils.InternalServerError("Error decoding request").WithInternalError(err)
		http_utils.SendErrorResponse(err, w, r)
		return
	}

	err = a.usvc.UpdatePassword(r.Context(), user, params.OldPassword, params.NewPassword)

	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}
	err = http_utils.SendJSON(w, http.StatusOK, &http_utils.StatusResponse{"ok"})
	if err != nil {
		log.Printf("Error sending response %v", err)
	}
}

type ResetPasswordParams struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (a *AuthApi) ResetPassword(w http.ResponseWriter, r *http.Request) {
	params := &ResetPasswordParams{}

	jsonDecoder := json.NewDecoder(r.Body)

	err := jsonDecoder.Decode(params)

	if err != nil {
		err = http_utils.InternalServerError("Error decoding request").WithInternalError(err)
		http_utils.SendErrorResponse(err, w, r)
	}

	err = a.usvc.ResetPassword(r.Context(), params.Token, params.Password)
	if err != nil {
		http_utils.SendErrorResponse(err, w, r)
		return
	}
	err = http_utils.SendJSON(w, http.StatusOK, &http_utils.StatusResponse{"ok"})
	if err != nil {
		log.Printf("Error sending response %v", err)
	}
}
